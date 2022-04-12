package ws_upgrader

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"homework15/level01-02/ws_parser"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// Headers' principle:

// Sec-WebSocket-Accept = encode_base64( sha1( Sec-WebSocket-Key + '258EAFA5-E914-47DA-95CA-C5AB0DC85B11' ))

// Resp: Sec-WebSocket-Protocol
// The server selects one or none of the acceptable protocols and echoes
// that value in its handshake to indicate that it has selected that
// protocol.

// 'If the client’s handshake did not contain such a header field(|Sec-WebSocket-Protocol|)
// or if the server does not agree to any of the client’s requested subprotocols, the
// only acceptable value is null. The absence of such a field is equivalent
// to the null value (meaning that if the server does not wish to agree to
// one of the suggested subprotocols, it MUST NOT send back a
// |Sec-WebSocket-Protocol| header field in its response).' —— RFC6455

// Example:

// HTTP/1.1 101 Switching Protocols
// Upgrade: websocket
// Connection: Upgrade
// Sec-WebSocket-Accept: GZS2YkUYBCu6eW6qTtiqD2bqEFE=

// GET / HTTP/1.1
// Host: 121.40.165.18:8800
// User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:99.0) Gecko/20100101 Firefox/99.0
// Accept: */*
// Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
// Accept-Encoding: gzip, deflate
// Sec-WebSocket-Version: 13
// Origin: http://www.websocket-test.com
// Sec-WebSocket-Extensions: permessage-deflate
// Sec-WebSocket-Key: HSohHIffGH1RFigAVUNDYw==
// Connection: keep-alive, Upgrade
// Pragma: no-cache
// Cache-Control: no-cache
// Upgrade: websocket

type WebsocketConn struct {
	mu         chan struct{} // protect write
	c          net.Conn
	pingHandle func(conn *WebsocketConn)
	Closed     bool
}

type Msg struct {
	Typ     ws_parser.OpenCodeType
	Content []byte
}

type Upgrader struct {
	PingHandle  func(conn *WebsocketConn)
	CheckOrigin func(origin string) bool
}

type UpgraderOption func(u *Upgrader)

var (
	HandShakeHeaderTemplate = "HTTP/1.1 101 Switching Protocols\nUpgrade: websocket\nConnection: Upgrade\nSec-WebSocket-Accept: {acc_key}\n\n"
	MagicHandShakeKey       = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	waitTimeout             = 10 * time.Second
	pongTimeout             = 60 * time.Second

	AllOriginOption UpgraderOption = func(u *Upgrader) {
		u.CheckOrigin = func(origin string) bool {
			return true
		}
	}

	DefaultPingHandle UpgraderOption = func(u *Upgrader) {
		u.PingHandle = func(conn *WebsocketConn) {
			frameBytes := ws_parser.WrapperPayload(nil, ws_parser.PongType)
			_, err := conn.c.Write(frameBytes)
			log.Println("err when response to PING, cause", err)
		}
	}
)

func NewUpgrader(options ...UpgraderOption) *Upgrader {
	u := &Upgrader{}
	for _, option := range options {
		option(u)
	}
	return u
}

func (u *Upgrader) Update(w http.ResponseWriter, r *http.Request) (conn *WebsocketConn, err error) {
	header := r.Header

	// detect if it is a normal request

	if !strings.Contains(header.Get("Connection"), "Upgrade") {
		return nil, nil
	}

	if header.Get("Upgrade") != "websocket" {
		return nil, nil
	}

	if r.Method != http.MethodGet {
		return nil, nil
	}

	if header.Get("Sec-Websocket-Version") != "13" {
		return nil, nil
	}

	// start hijack

	h, ok := w.(http.Hijacker)
	if !ok {
		err = errors.New("response writer cannot cast to hijacker")
		return
	}
	c, _, err := h.Hijack()
	if err != nil {
		err = errors.New("response writer cannot be hijacked")
		return
	}

	if !u.CheckOrigin(header.Get("Origin")) {
		err = c.Close() // close if origin cannot match
		return
	}

	rawKey := header.Get("Sec-Websocket-Key")

	if rawKey == "" {
		err = c.Close()
		return
	}

	conn = &WebsocketConn{
		mu:         make(chan struct{}, 1),
		c:          c,
		pingHandle: u.PingHandle,
	}
	conn.mu <- struct{}{}

	resp := []byte(strings.Replace(HandShakeHeaderTemplate, "{acc_key}", hashAccKey(rawKey), -1))
	_, err = c.Write(resp)
	if err != nil {
		err = errors.New("cannot write bytes to hijacked conn")
		return
	}

	// heart pump
	go func() {
		tick := time.NewTicker(pongTimeout)
		defer tick.Stop()
		for {
			err := conn.Ping()
			if err != nil {
				conn.Close()
				return
			}
			<-tick.C
		}
	}()

	return
}

func (c *WebsocketConn) Receive() (msg Msg, err error) {
	msg = Msg{}
	data := make([]byte, 1<<15)
	_, err = c.c.Read(data)
	if err != nil {
		c.c.Close() // regard as connect is close
		return
	}
	frame, err := ws_parser.ParseData(data)
	if err != nil {
		return
	}
	if !frame.FIN {
		err = errors.New("cannot handle FIN=0 condition")
		return
	}

	c.c.SetReadDeadline(time.Now().Add(pongTimeout))

	msg.Typ = frame.OpenCode

	switch frame.OpenCode {
	case ws_parser.PingType:
		c.pingHandle(c)
	case ws_parser.TextDataType, ws_parser.BinaryDataType:
		msg.Content = frame.Payload.GetData()
	case ws_parser.CloseConnType:
		c.Closed = true
		err = c.c.Close()
	case ws_parser.PongType:
		c.c.SetReadDeadline(time.Now().Add(pongTimeout))
	case ws_parser.ContinuationDataType:
		err = errors.New("cannot handle FIN=0 condition")
	}

	return
}

func (c *WebsocketConn) Close() error {
	c.Closed = true
	return c.c.Close()
}

func (c *WebsocketConn) Ping() error {
	<-c.mu
	defer func() { c.mu <- struct{}{} }()

	err := c.c.SetWriteDeadline(time.Now().Add(waitTimeout))
	if err != nil {
		return err
	}
	_, err = c.c.Write(ws_parser.WrapperPayload(nil, ws_parser.PingType))
	return err
}

func (c *WebsocketConn) Send(msg Msg) error {
	<-c.mu
	defer func() { c.mu <- struct{}{} }()

	if msg.Typ == ws_parser.CloseConnType {
		c.Closed = true
		defer c.c.Close()
	}
	err := c.c.SetWriteDeadline(time.Now().Add(waitTimeout))
	if err != nil {
		return err
	}
	_, err = c.c.Write(ws_parser.WrapperPayload(msg.Content, msg.Typ))
	return err
}

func hashAccKey(rawKey string) string {
	hash := sha1.New()
	hash.Write([]byte(rawKey + MagicHandShakeKey))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

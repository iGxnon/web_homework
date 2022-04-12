package chat_server

import (
	"homework15/level01-02/ws_parser"
	"homework15/level01-02/ws_upgrader"
	"regexp"
)

type Client struct {
	ID     string
	Conn   *ws_upgrader.WebsocketConn
	Send   chan []byte
	Closed bool

	Manager *ClientManager
}

func (c *Client) Close() error {
	c.Closed = true
	close(c.Send)
	return c.Conn.Close()
}

func (c *Client) ListenRead() {
	defer func() {
		c.Manager.UnRegister <- c
		c.Conn.Close()
	}()
	for {
		msg, err := c.Conn.Receive()
		if msg.Typ != ws_parser.TextDataType {
			continue
		}
		if err != nil || msg.Typ == ws_parser.CloseConnType {
			break
		}

		// simple way to locate id
		reg := regexp.MustCompile(`id\(\d\)`)
		find := reg.Find(msg.Content)
		tid := ""
		if len(find) != 0 {
			tid = string(find)[3 : len(find)-1]
		}

		mess := Message{
			SenderID:    c.ID,
			RecipientID: tid,
			Content:     string(msg.Content),
		}
		c.Manager.Broadcast <- mess
	}
}

func (c *Client) ListenWrite() {
	defer func() {
		c.Manager.UnRegister <- c
		c.Conn.Close()
	}()
	for {
		select {
		case mess := <-c.Send:
			err := c.Conn.Send(ws_upgrader.Msg{
				Typ:     ws_parser.TextDataType,
				Content: mess,
			})
			if err != nil {
				return
			}
		}
	}
}

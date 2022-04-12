package chat_server

import "homework15/level01-02/ws_upgrader"

type Client struct {
	ID     string
	Conn   *ws_upgrader.WebsocketConn
	Send   chan []byte
	Closed bool
}

func (c *Client) Close() error {
	c.Closed = true
	close(c.Send)
	return c.Conn.Close()
}

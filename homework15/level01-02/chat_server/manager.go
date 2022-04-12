package chat_server

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

type ClientManager struct {
	mu      *sync.RWMutex // keep write safe
	Clients map[*Client]struct{}

	ClientMap map[string]*Client

	Broadcast  chan Message
	Register   chan *Client
	UnRegister chan *Client
}

type Message struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id,omitempty"`
	Content     string `json:"content"`
}

func NewServer() *ClientManager {
	return &ClientManager{
		mu:        &sync.RWMutex{},
		Clients:   make(map[*Client]struct{}),
		ClientMap: make(map[string]*Client),

		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

func (c *ClientManager) ServeWithContext(ctx context.Context) {
	log.Println("Server start")
LOOP:
	for {
		select {
		case conn := <-c.Register:
			c.Clients[conn] = struct{}{}
			c.ClientMap[conn.ID] = conn

			mess := Message{
				SenderID: conn.ID,
				Content:  "came in",
			}
			log.Printf("A new client(%s) came in\n", conn.ID)
			c.BroadcastMessage(mess, conn)
		case conn := <-c.UnRegister:
			func() {
				c.mu.Lock()
				defer c.mu.Unlock()
				if _, ok := c.Clients[conn]; ok {
					delete(c.ClientMap, conn.ID)
					delete(c.Clients, conn)
					mess := Message{
						SenderID: conn.ID,
						Content:  "came out",
					}
					log.Printf("A client(%s) came out\n", conn.ID)
					c.BroadcastMessage(mess, conn)

					close(conn.Send)
				}
			}()
		case mess := <-c.Broadcast:
			log.Println("broadcast: ", mess)
			if mess.RecipientID == "" {
				c.BroadcastMessage(mess, c.ClientMap[mess.SenderID])
				continue
			}

			func() {
				c.mu.Lock()
				defer c.mu.Unlock()

				target, ok := c.ClientMap[mess.RecipientID]
				if !ok { // target not online
					sender := c.ClientMap[mess.SenderID]
					mess.Content = "target not online!"
					mess.RecipientID = mess.SenderID
					m, _ := json.Marshal(mess)
					sender.Send <- m
					return
				}
				m, _ := json.Marshal(mess)
				target.Send <- m
			}()
		case <-ctx.Done():
			func() {
				c.mu.Lock()
				defer c.mu.Unlock()
				for client := range c.Clients {
					if !client.Closed {
						err := client.Close()
						if err != nil {
							log.Printf("error when close client(%s) err: %v\n", client.ID, err)
						}
					}
				}

				c.Clients = nil
			}()

			close(c.Broadcast)
			close(c.Register)
			close(c.UnRegister)
			break LOOP
		}
	}
	log.Println("WebSocket Server has been closed...")
}

func (c *ClientManager) BroadcastBytes(mess []byte, ignore *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for client := range c.Clients {
		if client != ignore {
			client.Send <- mess
		}
	}
}

func (c *ClientManager) BroadcastMessage(mess Message, ignore *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	m, _ := json.Marshal(mess)
	for client := range c.Clients {
		if client != ignore {
			client.Send <- m
		}
	}
}

func (c *ClientManager) CheckOnline(client *Client) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Clients[client]
	return ok
}

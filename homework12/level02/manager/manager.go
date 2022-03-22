package manager

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

type Subscriber struct {
	ID         string
	conn       redis.Conn            // cannot create a new connection from other package
	Subscribed map[*Channel]struct{} // the channels this subscriber subscribed
}

func NewSubscriber(id string, conn redis.Conn) *Subscriber {
	return &Subscriber{
		ID:         id,
		conn:       conn,
		Subscribed: make(map[*Channel]struct{}),
	}
}

// Err returns a non-nil value when the connection is not usable.
func (s *Subscriber) Err() error {
	return s.conn.Err()
}

// Subscribe this method will block util this Channel handles
// the request
func (s *Subscriber) Subscribe(channel *Channel) {
	channel.Subscribe <- s
}

// Unsubscribe this method will block util this Channel handles
// the request
//
// When unsubscribe the channel,
func (s *Subscriber) Unsubscribe(channel *Channel) {
	channel.Unsubscribe <- s
}

func (s *Subscriber) Receive() interface{} {
	return redis.PubSubConn{Conn: s.conn}.Receive()
}

type Channel struct {
	ID          string
	Subscribers map[*Subscriber]struct{} // 订阅状态

	Publish     chan interface{} // broadcast message
	Subscribe   chan *Subscriber // 需要连接的
	Unsubscribe chan *Subscriber // 需要断开订阅的

	producer redis.Conn

	stopSignal chan struct{}
}

func NewChannel(id string) *Channel {
	return &Channel{
		ID:          id,
		Subscribers: make(map[*Subscriber]struct{}),
		Publish:     make(chan interface{}),
		Subscribe:   make(chan *Subscriber),
		Unsubscribe: make(chan *Subscriber),
		stopSignal:  make(chan struct{}),
	}
}

func (c *Channel) Serve(co redis.Conn) {
	c.producer = co
	for {
		select {
		case conn := <-c.Subscribe:
			// handle subscribe
			c.Subscribers[conn] = struct{}{}
			conn.Subscribed[c] = struct{}{}

			if err := conn.conn.Err(); err != nil {
				log.Printf("subsriber's connetction is not usable, err: %v\n", err)
				continue
			}
			_, err := conn.conn.Do("SUBSCRIBE", c.ID)
			if err != nil {
				log.Printf("Cannot handle this subscribe request from %s, err: %v\n", conn.ID, err)
				continue
			}
			log.Printf("A request from a subscriber %s subscribe this channel: %s, now subscriber count: %d\n", conn.ID, c.ID, len(c.Subscribers))
		case conn := <-c.Unsubscribe:
			if _, ok := c.Subscribers[conn]; ok {

				delete(c.Subscribers, conn)
				delete(conn.Subscribed, c)

				_, err := conn.conn.Do("UNSUBSCRIBE", c.ID)
				if err != nil {
					log.Printf("Cannot handle this unsubscribe request from %s, err: %v\n", conn.ID, err)
					continue
				}

				// todo remove?
				if len(conn.Subscribed) == 0 {
					err := conn.conn.Close()
					if err != nil {
						log.Printf("error occured when close a subscriber connection, err: %v\n", err)
						continue
					}
					log.Println("successfully close a subscriber's connection")
				}

				log.Printf("A request from a subscriber unsubscribe this channel: %s, now subscriber count: %d\n", c.ID, len(c.Subscribers))
			}
		case mess := <-c.Publish:
			_, err := c.producer.Do("PUBLISH", c.ID, mess)
			if err != nil {
				log.Printf("error occured when publish a mess, err: %v\n", err)
				continue
			}
			log.Println("A new mess had been sent")
		case _ = <-c.stopSignal:
			return
		}
	}
}

func (c *Channel) Close() error {
	if c.producer == nil {
		log.Panicf("cannot close a channel before Serve")
	}
	err := c.producer.Close()
	if err != nil {
		log.Printf("error occured when close connection, err: %v\n", err)
		return err
	}

	for subscriber := range c.Subscribers {
		_, err = subscriber.conn.Do("UNSUBSCRIBE", c.ID)
		if err != nil {
			log.Printf("Cannot unsubscribe from %s, err: %v\n", subscriber.ID, err)
			log.Println("Please release these subscribers by hand")
		}
		delete(subscriber.Subscribed, c)
	}

	c.stopSignal <- struct{}{}
	log.Printf("channel %s had been stopped", c.ID)
	return nil
}

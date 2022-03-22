package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
	"web_homework/homework12/level02/manager"
)

func main() {
	test()
}

func test() {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://root:@localhost:6379")
		},
		MaxIdle:   10,
		MaxActive: 15,
		Wait:      true,
	}

	defer pool.Close()

	c1 := manager.NewChannel("c1")

	s1 := manager.NewSubscriber("s1", pool.Get())
	s2 := manager.NewSubscriber("s2", pool.Get())

	go c1.Serve(pool.Get())

	// Subscribe invokes must after channel Serve
	s1.Subscribe(c1)
	s2.Subscribe(c1)

	<-time.After(time.Second)

	c1.Publish <- "Hi,there!"

	fmt.Println("message from c1 to s1: ", string(s1.Receive().(redis.Message).Data))
	fmt.Println("message from c1 to s2: ", string(s2.Receive().(redis.Message).Data))

	err := c1.Close()
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("s1 now subscribes %d channels\n", len(s1.Subscribed))
	fmt.Printf("s2 now subscribes %d channels\n", len(s2.Subscribed))
}

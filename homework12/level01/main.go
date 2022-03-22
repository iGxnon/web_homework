package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func main() {
	conn, err := redis.DialURL("redis://root:@localhost:6379")
	if err != nil {
		panic(err)
	}
	reply, err := conn.Do("PING")
	fmt.Println(reply, err)
	conn.Close()

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://root:@localhost:6379")
		},
		MaxIdle:   10,
		MaxActive: 15,
		Wait:      true,
	}
	conn = pool.Get()
	reply, err = conn.Do("SUBSCRIBE", "c1")
	replyList := reply.([]interface{})
	fmt.Println(string(replyList[0].([]byte)), string(replyList[1].([]byte)), replyList[2], err)
	go func() {
		for {
			receive, err := conn.Receive()
			receiveList := receive.([]interface{})
			fmt.Println(string(receiveList[0].([]byte)), string(receiveList[1].([]byte)), string(receiveList[2].([]byte)), err)
		}
	}()
	// 发布一条消息
	conn2 := pool.Get()
	reply, err = conn2.Do("PUBLISH", "c1", "hi, subscribers!")
	fmt.Println(reply, err)

	<-time.After(time.Second)
}

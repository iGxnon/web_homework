package main

import (
	"context"
	"homework15/level01-02/chat_server"
	"homework15/level01-02/ws_parser"
	"homework15/level01-02/ws_upgrader"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var id int64

func main() {
	chatRoom()
}

func chatRoom() {
	server := chat_server.NewServer()
	ctx, can := context.WithCancel(context.Background())
	defer can()
	go server.ServeWithContext(ctx)

	mux := http.NewServeMux()
	updater := ws_upgrader.NewUpgrader(ws_upgrader.DefaultPingHandle, ws_upgrader.AllOriginOption)
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		conn, err := updater.Update(w, r)
		log.Println(conn, err)
		client := &chat_server.Client{
			ID:      uid(),
			Conn:    conn,
			Send:    make(chan []byte),
			Manager: server,
		}
		server.Register <- client
		go client.ListenRead()
		go client.ListenWrite()
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func uid() string {
	atomic.AddInt64(&id, 1)
	return strconv.Itoa(int(id))
}

func test() {
	mux := http.NewServeMux()
	updater := ws_upgrader.NewUpgrader(ws_upgrader.DefaultPingHandle, ws_upgrader.AllOriginOption)
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		conn, err := updater.Update(w, r)
		log.Println(conn, err)
		// run a goroutine to listen request
		go func() {
			defer conn.Close()
			for !conn.Closed {
				msg, err2 := conn.Receive()
				if err2 != nil {
					log.Println(err2)
				}
				if msg.Typ == ws_parser.TextDataType {
					log.Println("Msg:", string(msg.Content))
				}
				if msg.Typ == ws_parser.CloseConnType {
					log.Println("A connect has been closed")
				}
			}
		}()
		// run a goroutine to send Hi!!! to the client
		tick := time.Tick(time.Second * 5)
		go func() {
			defer conn.Close()
			for !conn.Closed {
				err = conn.Send(ws_upgrader.Msg{
					Typ:     ws_parser.TextDataType,
					Content: []byte("Hi!!!"),
				})
				if err != nil {
					log.Println(err)
				}
				err = conn.Ping()
				if err != nil {
					log.Println(err)
				}
				<-tick
			}
		}()
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}

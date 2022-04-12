package main

import (
	"homework15/level01-02/ws_parser"
	"homework15/level01-02/ws_upgrader"
	"log"
	"net/http"
	"time"
)

func main() {
	test()
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

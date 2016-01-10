package main

import (
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/morikuni/go-training/chat/chat"
)

func main() {
	server := chat.NewServer()

	go server.Start()

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/chat", websocket.Handler(func(conn *websocket.Conn) {
		client := server.Client(conn)
		client.Start()
	}))

	http.ListenAndServe("localhost:9000", nil)
}

package chat

import (
	"log"

	"golang.org/x/net/websocket"
)

type Client interface {
	Start()
	Name() string
	Send(*MessageJSON)
	Message(string) *Message
	Close()
	SetRoom(Room)
}

func NewClient(conn *websocket.Conn, room Room) Client {
	c := &client{
		name: "unknown",
		conn: conn,
		room: room,
	}
	return c
}

type client struct {
	name string
	conn *websocket.Conn
	room Room
}

func (c *client) Name() string {
	return c.name
}

func (c *client) Message(msg string) *Message {
	return &Message{
		c,
		msg,
	}
}

func (c *client) Send(js *MessageJSON) {
	websocket.JSON.Send(c.conn, js)
}

func (c *client) Start() {
	var req ConnectionRequestJSON
	err := websocket.JSON.Receive(c.conn, &req)
	if err != nil {
		log.Println("request error ", err)
		c.Close()
		return
	}
	c.name = req.Name
	log.Println(c.name, "joined")

	for {
		var msg MessageJSON
		err := websocket.JSON.Receive(c.conn, &msg)
		if err != nil {
			log.Println(c.name, "exit by ", err)
			c.Close()
			return
		}
		c.room.Broadcast(&Message{c, msg.Message})
	}
}

func (c *client) Close() {
	c.room.Remove(c)
	c.conn.Close()
}

func (c *client) SetRoom(r Room) {
	c.room.Remove(c)
	c.room = r
}

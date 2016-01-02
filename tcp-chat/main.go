package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Message struct {
	client  *Client
	message []byte
}

type ChatRoom struct {
	clients []*Client
	c       chan []byte
	mu      sync.Mutex
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients: []*Client{},
		c:       make(chan []byte),
	}
}

func (r *ChatRoom) Start() {
	for {
		b := <-r.c
		for _, c := range r.clients {
			c.Send(b)
		}
	}
}

func (r *ChatRoom) Broadcast(msg Message) {
	r.c <- append(msg.client.prefix, msg.message...)
}

func (r *ChatRoom) Client(conn net.Conn) *Client {
	c := &Client{
		prefix: []byte(conn.RemoteAddr().String() + " : "),
		conn:   conn,
		room:   r,
	}
	r.mu.Lock()
	r.clients = append(r.clients, c)
	r.mu.Unlock()
	log.Println(conn.RemoteAddr().String(), "Joined")
	r.Broadcast(c.Message([]byte("JOINED\n")))
	return c
}

func (r *ChatRoom) Remove(c *Client) {
	r.mu.Lock()
	for i := range r.clients {
		if r.clients[i] == c {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			break
		}
	}
	r.mu.Unlock()
	log.Println(c.conn.RemoteAddr().String(), "Quit")
	r.Broadcast(c.Message([]byte("QUIT\n")))
	return
}

type Client struct {
	prefix []byte
	conn   net.Conn
	room   *ChatRoom
	quit   chan struct{}
}

func (c *Client) Send(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Start() {
	defer c.conn.Close()
	buf := make([]byte, 1024)
	for {
		c.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		s, err := c.conn.Read(buf)
		if err != nil {
			e := err
			if err, ok := err.(net.Error); ok && err.Timeout() {
				select {
				case <-c.quit:
					c.close()
					return
				default:
					continue
				}
			} else {
				c.close()
				log.Println(e)
				return
			}
		}
		c.room.Broadcast(c.Message(buf[:s]))
	}
}

func (c *Client) Message(b []byte) Message {
	return Message{
		client:  c,
		message: b,
	}
}

func (c *Client) Quit() {
	c.quit <- struct{}{}
}

func (c *Client) close() {
	c.room.Remove(c)
	c.conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
	}
	defer listener.Close()
	room := NewChatRoom()
	go room.Start()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}

		c := room.Client(conn)

		go c.Start()
	}
}

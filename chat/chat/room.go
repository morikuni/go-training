package chat

import (
	"sync"
)

type Room interface {
	Start()
	Name() string
	Add(Client)
	Remove(Client)
	Broadcast(*Message)
	Close()
}

func NewRoom(name string, s Server) Room {
	return &room{
		name:   name,
		server: s,
	}
}

type room struct {
	name    string
	server  Server
	clients []Client
	mu      sync.RWMutex
}

func (r *room) Start() {

}

func (r *room) Name() string {
	return r.name
}

func (r *room) Add(c Client) {
	r.mu.Lock()
	r.clients = append(r.clients, c)
	r.mu.Unlock()
}

func (r *room) Remove(c Client) {
	r.mu.RLock()
	idx := -1
	for i, cc := range r.clients {
		if c == cc {
			idx = i
			break
		}
	}
	r.mu.RUnlock()
	if idx == -1 {
		return
	}
	r.mu.Lock()
	r.clients = append(r.clients[:idx], r.clients[idx+1:]...)
	r.mu.Unlock()
}

func (r *room) Broadcast(m *Message) {
	js := &MessageJSON{
		m.Sender.Name(),
		m.Message,
	}
	r.mu.RLock()
	for _, c := range r.clients {
		c.Send(js)
	}
	r.mu.RUnlock()
}

func (r *room) Close() {
	r.server.Remove(r)
	r.mu.RLock()
	for _, c := range r.clients {
		c.Close()
	}
	r.mu.RUnlock()
}

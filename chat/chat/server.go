package chat

import (
	"sync"

	"golang.org/x/net/websocket"
)

type Server interface {
	Start()
	Room(name string) Room
	Remove(Room)
	Client(*websocket.Conn) Client
}

func NewServer() Server {
	s := &server{
		rooms: map[string]Room{},
	}
	s.lobby = NewRoom("lobby", s)
	return s
}

type server struct {
	rooms map[string]Room
	lobby Room
	mu    sync.RWMutex
}

func (s *server) Start() {

}

func (s *server) Room(name string) Room {
	s.mu.RLock()
	if room, ok := s.rooms[name]; ok {
		s.mu.RUnlock()
		return room
	}
	s.mu.Unlock()
	s.mu.Lock()
	room := NewRoom(name, s)
	s.rooms[name] = room
	s.mu.Unlock()
	return room
}

func (s *server) Remove(r Room) {
	s.mu.Lock()
	delete(s.rooms, r.Name())
	s.mu.Unlock()
}

func (s *server) Client(conn *websocket.Conn) Client {
	c := NewClient(conn, s.lobby)
	s.lobby.Add(c)
	return c
}

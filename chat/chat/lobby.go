package chat

func newLobby(s Server) *lobby {
	return &lobby{
		server:   s,
		delegate: NewRoom("lobby", s),
	}
}

type lobby struct {
	server   Server
	delegate Room
}

func (l *lobby) Start() {

}

func (l *lobby) Name() string {
	return "lobby"
}

func (l *lobby) Add(c Client) {
	l.delegate.Add(c)
}

func (l *lobby) Remove(c Client) {
	l.delegate.Remove(c)
}

func (l *lobby) Broadcast(m *Message) {
	roomName := m.Message
	if roomName == "lobby" {
		return
	}
	room := l.server.Room(roomName)
	room.Add(m.Sender)
	m.Sender.SetRoom(room)
}

func (l *lobby) Close() {
	l.delegate.Close()
}

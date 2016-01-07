package chat

type Message struct {
	Sender  Client
	Message string
}

type MessageJSON struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type ConnectionRequestJSON struct {
	Name string `json:"name"`
}

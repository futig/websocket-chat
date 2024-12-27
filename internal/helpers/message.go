package helpers

type Message struct {
	MType string `json:"mtype"`
	ID    string `json:"id"`
	To    string `json:"to,omitempty"`
	Text  string `json:"text,omitempty"`
}

type DirectMessage struct {
	To      string
	Message Message
}

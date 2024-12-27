package interfaces

import (
	h "github.com/futig/websocket-chat/internal/helpers"
)

type Client interface {
	ID() string
	Send() chan h.Message
}

type Hub interface {
	Register(client Client)
	Unregister(client Client)
	Broadcast(message h.Message)
	SendDirect(toID string, message h.Message)
}

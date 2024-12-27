package hub

import (
	hp "github.com/futig/websocket-chat/internal/helpers"
	ifaces "github.com/futig/websocket-chat/internal/interfaces"
	"sync"
)

type Hub struct {
	clients    map[string]ifaces.Client

	register   chan ifaces.Client
	unregister chan ifaces.Client
	broadcast  chan hp.Message
	direct     chan hp.DirectMessage

	mu         sync.RWMutex
}

func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[string]ifaces.Client),
		register:   make(chan ifaces.Client),
		unregister: make(chan ifaces.Client),
		broadcast:  make(chan hp.Message),
		direct:     make(chan hp.DirectMessage),
	}
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID()] = client
			h.mu.Unlock()
			enterMsg := hp.Message{
				MType: "USER_ENTER",
				ID:    client.ID(),
			}
			h.Broadcast(enterMsg)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID()]; ok {
				delete(h.clients, client.ID())
				close(client.Send())
				leaveMsg := hp.Message{
					MType: "USER_LEAVE",
					ID:    client.ID(),
				}
				h.Broadcast(leaveMsg)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send() <- message:
				default:
					close(client.Send())
					delete(h.clients, client.ID())
				}
			}
			h.mu.RUnlock()
		case direct := <-h.direct:
			h.mu.RLock()
			if client, ok := h.clients[direct.To]; ok {
				select {
				case client.Send() <- direct.Message:
				default:
					close(client.Send())
					delete(h.clients, client.ID())
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Register(client ifaces.Client) {
	h.register <- client
}

func (h *Hub) Unregister(client ifaces.Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(message hp.Message) {
	h.broadcast <- message
}

func (h *Hub) SendDirect(toID string, message hp.Message) {
	h.direct <- hp.DirectMessage{
		To:      toID,
		Message: message,
	}
}

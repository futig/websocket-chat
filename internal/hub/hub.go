package hub

import (
	"sync"

	hp "github.com/futig/websocket-chat/internal/helpers"
	ifaces "github.com/futig/websocket-chat/internal/interfaces"
)

type Hub struct {
	clients map[string]ifaces.Client

	register   chan ifaces.Client
	unregister chan ifaces.Client
	broadcast  chan hp.Message
	direct     chan hp.DirectMessage

	mu sync.RWMutex
}

func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[string]ifaces.Client),
		register:   make(chan ifaces.Client, 100),
		unregister: make(chan ifaces.Client, 100),
		broadcast:  make(chan hp.Message, 100),
		direct:     make(chan hp.DirectMessage, 100),
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
			ignore := ""
			if message.MType == "MSG" {
				ignore = message.ID
			}
			h.mu.RLock()
			for _, client := range h.clients {
				if client.ID() != ignore {
					select {
					case client.Send() <- message:
					default:
						close(client.Send())
						delete(h.clients, client.ID())
					}
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
		default:
			continue
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

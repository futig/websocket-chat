package client

import (
	"encoding/json"
	"log"
	"time"

	h "github.com/futig/websocket-chat/internal/helpers"
	ifaces "github.com/futig/websocket-chat/internal/interfaces"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	conn   *websocket.Conn
	send   chan h.Message
	hub    ifaces.Hub
	closed bool
}

func NewClient(id string, conn *websocket.Conn, hub ifaces.Hub) *Client {
	return &Client{
		id:   id,
		conn: conn,
		send: make(chan h.Message, 256),
		hub:  hub,
	}
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Send() chan h.Message {
	return c.send
}

func (c *Client) close() {
	if !c.closed {
		c.conn.Close()
		c.closed = true
	}
}

func (c *Client) ReadPump() {
	defer func() {
        c.hub.Unregister(c)
        c.close()
    }()

	c.conn.SetReadLimit(512)
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

	for {
		_, msgData, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		if string(msgData) == "ping" {
			c.conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}

		var msg h.Message
		err = json.Unmarshal(msgData, &msg)
		if err != nil {
			log.Println("invalid message format:", err)
			continue
		}

		switch msg.MType {
		case "TEXT":
			if msg.To == "" {
				outMsg := h.Message{
					MType: "MSG",
					ID:    c.id,
					Text:  msg.Text,
				}
				c.hub.Broadcast(outMsg)
			} else {
				outMsg := h.Message{
					MType: "DM",
					ID:    c.id,
					Text:  msg.Text,
				}
				c.hub.SendDirect(msg.To, outMsg)
			}
		default:
			log.Println("unknown message type:", msg.MType)
		}
	}
}

func (c *Client) WritePump() {
	defer c.close()
	for msg := range c.send {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Println("marshal error:", err)
			continue
		}
		err = c.conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/futig/websocket-chat/internal/client"
	hp "github.com/futig/websocket-chat/internal/helpers"
	"github.com/futig/websocket-chat/internal/hub"

	"github.com/gorilla/websocket"
)

type Server struct {
	Port   uint16
	IpAddr string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewServer(ip string, port uint16) *Server {
	return &Server{
		Port:   port,
		IpAddr: ip,
	}
}

func serveWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	_, initMsgData, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}

	var initMsg hp.Message
	err = json.Unmarshal(initMsgData, &initMsg)
	if err != nil || initMsg.MType != "INIT" || len(initMsg.ID) == 0 {
		conn.Close()
		return
	}

	client := client.NewClient(initMsg.ID, conn, hub)
	hub.Register(client)

	go client.ReadPump()
	go client.WritePump()
}

// serveHome отдаёт index.html
func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func (s *Server) Run() {
	hub := hub.NewHub()
	go hub.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	addr := fmt.Sprintf("%s:%d", s.IpAddr, s.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

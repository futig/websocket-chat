package main

import (
	"github.com/futig/websocket-chat/internal/server"
)

func main() {
    wsServer := server.NewServer("0.0.0.0", 8080)
	wsServer.Run()
}
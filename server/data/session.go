package data

import "golang.org/x/net/websocket"

type Session struct {
	connections map[*websocket.Conn]bool
}

func NewSession() *Session {
	return &Session{
		connections: make(map[*websocket.Conn]bool),
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/CoreCreation/scrum-poker/server/data"
	"github.com/CoreCreation/scrum-poker/server/handlers"
)

// func (s *Server) handleWS(ws *websocket.Conn) {
// 	fmt.Println("New incoming connection from client:", ws.RemoteAddr())

// 	// Need mutex? Not concurrent safe
// 	s.conns[ws] = true

// 	s.readLoop(ws)
// }

// func (s *Server) readLoop(ws *websocket.Conn) {
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := ws.Read(buf)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			fmt.Println("read error:", err)
// 			continue
// 		}
// 		msg := buf[:n]

// 		s.broadcast(msg)
// 	}
// }

// func (s *Server) broadcast(b []byte) {
// 	for ws := range s.conns {
// 		go func(ws *websocket.Conn) {
// 			if _, err := ws.Write(b); err != nil {
// 				fmt.Println("Write error:", err)
// 			}
// 		}(ws)
// 	}
// }

func main() {
	data := data.NewData()
	const addr = ":3001"

	// Handlers
	handlers := handlers.NewHandlers(data)
	router := handlers.GetRouter()

	// SPA Handler
	router.Handle("/", http.FileServer(http.Dir("dist")))

	// WebSocket Handler
	// http.Handle("/ws", websocket.Handler(server.handleWS))

	// Blocking loop
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

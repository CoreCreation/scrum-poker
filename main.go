package main

import (
	"log"
	"net/http"

	"github.com/CoreCreation/scrum-poker/server/data"
	"github.com/CoreCreation/scrum-poker/server/handlers"
)

func main() {
	data := data.NewData()
	const addr = ":3001"

	// Handlers
	handlers := handlers.NewHandlers(data)
	router := handlers.GetRouter()

	// SPA Handler
	router.Handle("/", http.FileServer(http.Dir("dist")))

	// Blocking loop
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CoreCreation/scrum-poker/server/data"
	"github.com/CoreCreation/scrum-poker/server/handlers"
)

func main() {
	fs := http.Dir("./dist")
	data := data.NewSessions()
	const addr = ":3001"

	fmt.Println("Starting Server at port:", addr)

	// Handlers
	handlers := handlers.NewHandlers(data)
	router := handlers.GetRouter()

	// SPA Handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := fs.Open(r.URL.Path)
		if err != nil {
			http.ServeFile(w, r, "./dist/index.html")
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			http.ServeFile(w, r, "./dist/index.html")
			return
		}

		http.FileServer(fs).ServeHTTP(w, r)
	})

	// Blocking loop
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

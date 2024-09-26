package main

import (
	"log"
	"net/http"
	// "sync"

	"SharkLava/random_chat/internal/handlers"
	"SharkLava/random_chat/pkg/queue"
)

func main() {
	userQueue := queue.NewQueue()
	hub := handlers.NewHub(userQueue)
	go hub.Run()

	http.HandleFunc("/", handlers.ServeHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(hub, w, r)
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

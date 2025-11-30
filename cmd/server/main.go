package main

import (
	"go_server/internal/config"
	"go_server/internal/database"
	"go_server/internal/handlers"
	"log"
	"net/http"
)

func main() {
	database.Init(config.DatabaseURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/status", handlers.HandleGetStatus)
	mux.HandleFunc("/users", handlers.HandleGetUsers)
	mux.HandleFunc("/create_event", handlers.HandleCreateEvent)

	log.Println("Server is running at :8080")
	http.ListenAndServe(":8080", mux)
}

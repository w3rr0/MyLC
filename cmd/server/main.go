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
	mux.HandleFunc("/delete_event", handlers.HandleDeleteEvent)
	mux.HandleFunc("/change_availability", handlers.HandleChangeAvailability)

	log.Println("Server is running at :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

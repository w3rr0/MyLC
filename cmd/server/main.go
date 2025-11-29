package main

import (
	"go_server/internal/config"
	"go_server/internal/database"
	"go_server/internal/handlers"
	"io"
	"log"
	"net/http"
)

func main() {
	database.Init(config.DatabaseURL)

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Ready")
	})
	http.HandleFunc("/users", handlers.GetUsers)

	log.Println("Server is running at :8080")
	http.ListenAndServe(":8080", nil)
}

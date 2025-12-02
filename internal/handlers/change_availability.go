package handlers

import (
	"encoding/json"
	"go_server/internal/database"
	"go_server/internal/repository"
	"net/http"
)

type RequestChangeAvailability struct {
	EventId      int               `json:"event_id"`
	UserId       int               `json:"user_id"`
	Availability map[string]string `json:"availability"`
}

func HandleChangeAvailability(w http.ResponseWriter, r *http.Request) {
	db := database.DB
	var req RequestChangeAvailability
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.ChangeAvailability(db, req.EventId, req.UserId, req.Availability)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

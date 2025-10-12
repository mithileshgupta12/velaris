package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/db"
)

type PostHandler struct {
	database *db.DB
}

func NewPostHandler(database *db.DB) *PostHandler {
	return &PostHandler{database}
}

func (ph *PostHandler) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := ph.database.Queries.GetAllPosts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

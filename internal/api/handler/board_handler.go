package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/db"
)

type BoardHandler struct {
	database *db.DB
}

func NewBoardHandler(database *db.DB) *BoardHandler {
	return &BoardHandler{database}
}

func (bh *BoardHandler) Index(w http.ResponseWriter, r *http.Request) {
	boards, err := bh.database.Queries.GetAllBoards(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(boards); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (bh *BoardHandler) Store(w http.ResponseWriter, r *http.Request) {
	//
}

func (bh *BoardHandler) Show(w http.ResponseWriter, r *http.Request) {
	//
}

func (bh *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	//
}

func (bh *BoardHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	//
}

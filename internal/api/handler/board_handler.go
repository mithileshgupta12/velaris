package handler

import (
	"log"
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/helper"
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
		log.Printf("failed to get boards: %v", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, boards)
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

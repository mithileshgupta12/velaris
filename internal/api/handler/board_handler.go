package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

type CreateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BoardHandler struct {
	database *db.DB
	lgr      logger.Logger
}

func NewBoardHandler(database *db.DB, lgr logger.Logger) *BoardHandler {
	return &BoardHandler{database, lgr}
}

func (bh *BoardHandler) Index(w http.ResponseWriter, r *http.Request) {
	boards, err := bh.database.Queries.GetAllBoards(r.Context())
	if err != nil {
		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to get boards: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, boards)
}

func (bh *BoardHandler) Store(w http.ResponseWriter, r *http.Request) {
	var createBoardRequest CreateBoardRequest

	if err := json.NewDecoder(r.Body).Decode(&createBoardRequest); err != nil {
		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to decode request: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	createBoardRequest.Name = strings.TrimSpace(createBoardRequest.Name)
	createBoardRequest.Description = strings.TrimSpace(createBoardRequest.Description)

	if createBoardRequest.Name == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name is a required field")
		return
	}

	createBoardParams := repository.CreateBoardParams{
		Name: createBoardRequest.Name,
	}

	if createBoardRequest.Description == "" {
		createBoardParams.Description = pgtype.Text{
			String: "",
			Valid:  false,
		}
	} else {
		createBoardParams.Description = pgtype.Text{
			String: createBoardRequest.Description,
			Valid:  true,
		}
	}

	board, err := bh.database.Queries.CreateBoard(r.Context(), createBoardParams)
	if err != nil {
		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to create board: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, board)
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

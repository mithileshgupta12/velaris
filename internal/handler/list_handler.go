package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
)

type ListHandler struct {
	listRepository repository.ListRepository
}

func NewListHandler(listRepository repository.ListRepository) *ListHandler {
	return &ListHandler{listRepository}
}

func (lh *ListHandler) Index(w http.ResponseWriter, r *http.Request) {
	boardIdParam := chi.URLParam(r, "boardId")

	boardId, err := strconv.Atoi(boardIdParam)
	if err != nil {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid board id")
		return
	}

	lists, err := lh.listRepository.GetAllListsByBoardId(&repository.GetAllListsByBoardIdArgs{
		BoardId: int64(boardId),
	})
	if err != nil {
		slog.Error("failed to get lists for board", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, lists)
}

func (lh *ListHandler) Store(w http.ResponseWriter, r *http.Request) {
	//
}

func (lh *ListHandler) Show(w http.ResponseWriter, r *http.Request) {
	//
}

func (lh *ListHandler) Update(w http.ResponseWriter, r *http.Request) {
	//
}

func (lh *ListHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	//
}

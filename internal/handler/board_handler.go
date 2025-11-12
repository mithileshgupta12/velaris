package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/db/policy"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/middleware"
)

type CreateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BoardHandler struct {
	boardRepository repository.BoardRepository
	boardPolicy     policy.Policy
}

func NewBoardHandler(boardRepository repository.BoardRepository, boardPolicy policy.Policy) *BoardHandler {
	return &BoardHandler{boardRepository, boardPolicy}
}

func (bh *BoardHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	boards, err := bh.boardRepository.GetAllBoardsByUserId(ctxUser.ID)
	if err != nil {
		slog.Error("failed to get boards", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, boards)
}

func (bh *BoardHandler) Store(w http.ResponseWriter, r *http.Request) {
	var createBoardRequest CreateBoardRequest

	if err := json.NewDecoder(r.Body).Decode(&createBoardRequest); err != nil {
		slog.Error("failed to decode request", "err", err)
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	createBoardRequest.Name = strings.TrimSpace(createBoardRequest.Name)
	createBoardRequest.Description = strings.TrimSpace(createBoardRequest.Description)

	if createBoardRequest.Name == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name is a required field")
		return
	}

	if len(createBoardRequest.Name) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name must not be more than 255 characters long")
		return
	}

	if len(createBoardRequest.Description) > 10000 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "description must not be more than 10,000 characters long")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	createBoardArgs := &repository.CreateBoardArgs{
		Name:   createBoardRequest.Name,
		UserId: ctxUser.ID,
	}

	if createBoardRequest.Description == "" {
		createBoardArgs.Description = nil
	} else {
		createBoardArgs.Description = &createBoardRequest.Description
	}

	board, err := bh.boardRepository.CreateBoard(createBoardArgs)
	if err != nil {
		slog.Error("failed to create board", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusCreated, board)
}

func (bh *BoardHandler) Show(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid board id")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	canView, err := bh.boardPolicy.CanView(ctxUser, int64(id))
	if err != nil {
		slog.Error("failed to check board view permission", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !canView {
		helper.ErrorJsonResponse(w, http.StatusForbidden, "unauthorized")
		return
	}

	board, err := bh.boardRepository.GetBoardById(&repository.GetBoardByIdArgs{
		Id: int64(id),
	})
	if err != nil {
		slog.Error("failed to get board by ID", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, board)
}

func (bh *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid board id")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	canUpdate, err := bh.boardPolicy.CanUpdate(ctxUser, int64(id))
	if err != nil {
		slog.Error("failed to check board update permission", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !canUpdate {
		helper.ErrorJsonResponse(w, http.StatusForbidden, "unauthorized")
		return
	}

	var updateBoardRequest UpdateBoardRequest

	if err := json.NewDecoder(r.Body).Decode(&updateBoardRequest); err != nil {
		slog.Error("failed to decode request", "err", err)
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	updateBoardRequest.Name = strings.TrimSpace(updateBoardRequest.Name)
	updateBoardRequest.Description = strings.TrimSpace(updateBoardRequest.Description)

	if updateBoardRequest.Name == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name is a required field")
		return
	}

	if len(updateBoardRequest.Name) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name must not be more than 255 characters long")
		return
	}

	if len(updateBoardRequest.Description) > 10000 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "description must not be more than 10,000 characters long")
		return
	}

	updateBoardByIdArgs := &repository.UpdateBoardByIdArgs{
		Id:   int64(id),
		Name: updateBoardRequest.Name,
	}

	if updateBoardRequest.Description == "" {
		updateBoardByIdArgs.Description = nil
	} else {
		updateBoardByIdArgs.Description = &updateBoardRequest.Description
	}

	board, err := bh.boardRepository.UpdateBoardById(updateBoardByIdArgs)
	if err != nil {
		slog.Error("failed to update board", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, board)
}

func (bh *BoardHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid board id")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	canDelete, err := bh.boardPolicy.CanDelete(ctxUser, int64(id))
	if err != nil {
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !canDelete {
		helper.ErrorJsonResponse(w, http.StatusForbidden, "unauthorized")
		return
	}

	err = bh.boardRepository.DeleteBoardById(&repository.DeleteBoardByIdArgs{
		Id: int64(id),
	})
	if err != nil {
		slog.Error("failed to delete board", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

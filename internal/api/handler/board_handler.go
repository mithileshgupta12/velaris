package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
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
	queries repository.Querier
	lgr     logger.Logger
}

func NewBoardHandler(queries repository.Querier, lgr logger.Logger) *BoardHandler {
	return &BoardHandler{queries, lgr}
}

func (bh *BoardHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	boards, err := bh.queries.GetAllBoardsByUserId(r.Context(), ctxUser.ID)
	if err != nil {
		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to get boards: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
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

	if len(createBoardRequest.Name) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name must not be more than 255 characters long")
		return
	}

	if len(createBoardRequest.Description) > 10000 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "description must not be more than 10,000 characters long")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	createBoardParams := repository.CreateBoardParams{
		Name:   createBoardRequest.Name,
		UserID: ctxUser.ID,
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

	board, err := bh.queries.CreateBoard(r.Context(), createBoardParams)
	if err != nil {
		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to create board: %v", err), nil)
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

	board, err := bh.queries.GetBoardByIdAndUserId(r.Context(), repository.GetBoardByIdAndUserIdParams{
		ID:     int64(id),
		UserID: ctxUser.ID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			helper.ErrorJsonResponse(w, http.StatusNotFound, "board not found")
			return
		}

		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to get board by ID: %v", err), []*logger.Field{
			{Key: "board_id", Value: id},
		})
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

	var updateBoardRequest UpdateBoardRequest

	if err := json.NewDecoder(r.Body).Decode(&updateBoardRequest); err != nil {
		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to decode request: %v", err), nil)
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

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	updateBoardByIdParams := repository.UpdateBoardByIdAndUserIdParams{
		ID:     int64(id),
		UserID: ctxUser.ID,
		Name:   updateBoardRequest.Name,
	}

	if updateBoardRequest.Description == "" {
		updateBoardByIdParams.Description = pgtype.Text{
			String: "",
			Valid:  false,
		}
	} else {
		updateBoardByIdParams.Description = pgtype.Text{
			String: updateBoardRequest.Description,
			Valid:  true,
		}
	}

	board, err := bh.queries.UpdateBoardByIdAndUserId(r.Context(), updateBoardByIdParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			helper.ErrorJsonResponse(w, http.StatusNotFound, "board not found")
			return
		}

		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to update board: %v", err), []*logger.Field{
			{Key: "board_id", Value: id},
		})
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

	rowsAffected, err := bh.queries.DeleteBoardByIdAndUserId(r.Context(), repository.DeleteBoardByIdAndUserIdParams{
		ID:     int64(id),
		UserID: ctxUser.ID,
	})
	if err != nil {
		bh.lgr.Log(logger.ERROR, fmt.Sprintf("failed to delete board: %v", err), []*logger.Field{
			{Key: "board_id", Value: id},
		})
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if rowsAffected == 0 {
		helper.ErrorJsonResponse(w, http.StatusNotFound, "board not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

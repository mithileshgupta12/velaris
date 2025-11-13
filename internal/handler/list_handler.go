package handler

import (
	"log/slog"
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/db/policy"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/middleware"
)

type ListHandler struct {
	listRepository repository.ListRepository
	boardPolicy    policy.Policy
	listPolicy     policy.Policy
}

func NewListHandler(
	listRepository repository.ListRepository,
	boardPolicy policy.Policy,
	listPolicy policy.Policy,
) *ListHandler {
	return &ListHandler{listRepository, boardPolicy, listPolicy}
}

func (lh *ListHandler) Index(w http.ResponseWriter, r *http.Request) {
	boardId, err := helper.ParseIntURLParam(r, "boardId")
	if err != nil || boardId < 1 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid board id")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	canView, err := lh.boardPolicy.CanView(ctxUser, boardId)
	if err != nil {
		slog.Error("failed to check board view permission", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !canView {
		helper.ErrorJsonResponse(w, http.StatusNotFound, "board not found")
		return
	}

	lists, err := lh.listRepository.GetAllListsByBoardId(&repository.GetAllListsByBoardIdArgs{
		BoardId: boardId,
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

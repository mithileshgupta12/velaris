package handler

import (
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type ListHandler struct {
	listRepository repository.ListRepository
}

func NewListHandler(listRepository repository.ListRepository) *ListHandler {
	return &ListHandler{listRepository}
}

func (lh *ListHandler) Index(w http.ResponseWriter, r *http.Request) {
	//
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

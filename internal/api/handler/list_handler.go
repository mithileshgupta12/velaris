package handler

import (
	"net/http"
)

type ListHandler struct {
}

func NewListHandler() *ListHandler {
	return &ListHandler{}
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

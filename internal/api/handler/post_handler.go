package handler

import (
	"log"
	"net/http"
)

type PostHandler struct{}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (ph *PostHandler) Index(w http.ResponseWriter, r *http.Request) {
	log.Println("Posts!")
}

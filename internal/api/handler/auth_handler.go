package handler

import (
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

type AuthHandler struct {
	queries repository.Querier
	lgr     logger.Logger
}

func NewAuthHandler(queries repository.Querier, lgr logger.Logger) *AuthHandler {
	return &AuthHandler{queries, lgr}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	//
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	//
}

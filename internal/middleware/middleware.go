package middleware

import (
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type Middlewares interface {
	AuthMiddleware(next http.Handler) http.Handler
}

type middlewares struct {
	repositories *repository.Repository
	sessionStore cache.SessionStore
}

func NewMiddlewares(repositories *repository.Repository, sessionStore cache.SessionStore) Middlewares {
	return &middlewares{repositories, sessionStore}
}

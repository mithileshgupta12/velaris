package middleware

import (
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

type Middlewares interface {
	AuthMiddleware(next http.Handler) http.Handler
}

type middlewares struct {
	lgr          logger.Logger
	repositories *repository.Repository
	sessionStore cache.SessionStore
}

func NewMiddlewares(lgr logger.Logger, repositories *repository.Repository, sessionStore cache.SessionStore) Middlewares {
	return &middlewares{lgr, repositories, sessionStore}
}

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
	queries      repository.Querier
	sessionStore cache.SessionStore
}

func NewMiddlewares(lgr logger.Logger, queries repository.Querier, sessionStore cache.SessionStore) Middlewares {
	return &middlewares{lgr, queries, sessionStore}
}

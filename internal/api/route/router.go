package route

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

type Router struct {
	mux *chi.Mux
	lgr logger.Logger
}

func NewRouter(lgr logger.Logger) *Router {
	mux := chi.NewRouter()

	mux.Use(chiMiddlewares.RequestID)
	mux.Use(chiMiddlewares.RealIP)
	mux.Use(chiMiddlewares.Logger)
	mux.Use(chiMiddlewares.Recoverer)

	return &Router{mux, lgr}
}

func (r *Router) RegisterRoutes(queries repository.Querier, stores *cache.Stores, middlewares middleware.Middlewares) {
	BoardRoutes(r.mux, queries, r.lgr)
	AuthRoutes(r.mux, queries, stores.SessionStore, r.lgr, middlewares)
}

func (r *Router) Serve(port int) error {
	addr := fmt.Sprintf(":%d", port)

	r.lgr.Log(logger.INFO, fmt.Sprintf("Server started on %s", addr), nil)
	return http.ListenAndServe(addr, r.mux)
}

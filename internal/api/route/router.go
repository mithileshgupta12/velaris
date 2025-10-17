package route

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

type Router struct {
	lgr      logger.Logger
	mux      *chi.Mux
	database *db.DB
}

func NewRouter(lgr logger.Logger, database *db.DB) *Router {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return &Router{lgr, mux, database}
}

func (r *Router) RegisterRoutes() {
	BoardRoutes(r.mux, r.database)
}

func (r *Router) Serve(port int) error {
	addr := fmt.Sprintf(":%d", port)

	r.lgr.Log(logger.FormatJSON, logger.INFO, fmt.Sprintf("Server started on %s", addr))
	return http.ListenAndServe(addr, r.mux)
}

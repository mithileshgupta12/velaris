package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mithileshgupta12/velaris/internal/db"
)

type Router struct {
	mux      *chi.Mux
	database *db.DB
}

func NewRouter(database *db.DB) *Router {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return &Router{mux, database}
}

func (r *Router) RegisterRoutes() {
	PostRoutes(r.mux, r.database)
}

func (r *Router) Serve(port int) error {
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Server started on %s", addr)
	return http.ListenAndServe(addr, r.mux)
}

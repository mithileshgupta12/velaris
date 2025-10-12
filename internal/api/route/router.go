package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type Router struct {
	mux     *chi.Mux
	queries *repository.Queries
}

func NewRouter(queries *repository.Queries) *Router {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return &Router{mux, queries}
}

func (r *Router) RegisterRoutes() {
	PostRoutes(r.mux)
}

func (r *Router) Serve(port int) error {
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Server started on %s", addr)
	return http.ListenAndServe(addr, r.mux)
}

package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	r *chi.Mux
}

func NewRouter() *Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	PostRoutes(r)

	return &Router{
		r,
	}
}

func (r *Router) Serve(port int) error {
	log.Printf("Server started on http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r.r); err != nil {
		return err
	}

	return nil
}

package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		mux: chi.NewRouter(),
	}
}

func (r *Router) Init() *Router {
	r.mux.Use(middleware.RequestID)
	r.mux.Use(middleware.RealIP)
	r.mux.Use(middleware.Logger)
	r.mux.Use(middleware.Recoverer)

	r.mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	r.mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	return r
}

func (r *Router) Serve(addr string) error {
	if err := http.ListenAndServe(addr, r.mux); err != nil {
		return err
	}

	return nil
}

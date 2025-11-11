package route

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter(frontendUrl string) *Router {
	mux := chi.NewRouter()

	mux.Use(chiMiddlewares.RequestID)
	mux.Use(chiMiddlewares.RealIP)
	mux.Use(chiMiddlewares.Logger)
	mux.Use(chiMiddlewares.Recoverer)
	mux.Use(middleware.LimitBodySize(1024 * 1024))

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return &Router{mux}
}

func (r *Router) RegisterRoutes(repositories *repository.Repository, stores *cache.Stores, middlewares middleware.Middlewares) {
	BoardRoutes(r.mux, repositories.BoardRepository, middlewares)
	AuthRoutes(r.mux, repositories.UserRepository, stores.SessionStore, middlewares)
}

func (r *Router) Serve(port int) error {
	addr := fmt.Sprintf(":%d", port)

	slog.Info(fmt.Sprintf("Server started on %s", addr))
	return http.ListenAndServe(addr, r.mux)
}

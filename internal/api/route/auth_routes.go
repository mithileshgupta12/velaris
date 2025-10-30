package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func AuthRoutes(r *chi.Mux, queries repository.Querier, sessionStore cache.SessionStore, lgr logger.Logger) {
	authHandler := handler.NewAuthHandler(queries, sessionStore, lgr)

	protected := r.With(middleware.AuthMiddleware)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		protected.Post("/logout", authHandler.Logout)
	})
}

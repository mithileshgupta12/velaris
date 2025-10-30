package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func AuthRoutes(
	r *chi.Mux,
	queries repository.Querier,
	sessionStore cache.SessionStore,
	lgr logger.Logger,
	middlewares middleware.Middlewares,
) {
	authHandler := handler.NewAuthHandler(queries, sessionStore, lgr)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)
			r.Post("/logout", authHandler.Logout)
			r.Get("/user", authHandler.GetLoggedInUser)
		})
	})
}

package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/handler"
	"github.com/mithileshgupta12/velaris/internal/middleware"
)

func AuthRoutes(
	r *chi.Mux,
	userRepository repository.UserRepository,
	sessionStore cache.SessionStore,
	middlewares middleware.Middlewares,
) {
	authHandler := handler.NewAuthHandler(userRepository, sessionStore)

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

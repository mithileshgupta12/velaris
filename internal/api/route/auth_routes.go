package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func AuthRoutes(r *chi.Mux, queries repository.Querier, lgr logger.Logger) {
	authHandler := handler.NewAuthHandler(queries, lgr)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})
}

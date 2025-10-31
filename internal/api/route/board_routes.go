package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func BoardRoutes(r *chi.Mux, queries repository.Querier, lgr logger.Logger, middlewares middleware.Middlewares) {
	boardHandler := handler.NewBoardHandler(queries, lgr)

	r.Route("/boards", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/", boardHandler.Index)
		r.Post("/", boardHandler.Store)
		r.Get("/{id}", boardHandler.Show)
		r.Put("/{id}", boardHandler.Update)
		r.Delete("/{id}", boardHandler.Destroy)
	})

}

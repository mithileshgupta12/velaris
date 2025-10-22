package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func BoardRoutes(r *chi.Mux, queries repository.Querier, lgr logger.Logger) {
	boardHandler := handler.NewBoardHandler(queries, lgr)

	r.Get("/boards", boardHandler.Index)
	r.Post("/boards", boardHandler.Store)
	r.Get("/boards/{id}", boardHandler.Show)
	r.Put("/boards/{id}", boardHandler.Update)
	r.Delete("/boards/{id}", boardHandler.Destroy)
}

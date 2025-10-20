package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func BoardRoutes(r *chi.Mux, database *db.DB, lgr logger.Logger) {
	boardHandler := handler.NewBoardHandler(database, lgr)

	r.Get("/boards", boardHandler.Index)
	r.Post("/boards", boardHandler.Store)
	r.Get("/boards/{id}", boardHandler.Store)
	r.Put("/boards/{id}", boardHandler.Store)
	r.Delete("/boards/{id}", boardHandler.Store)
}

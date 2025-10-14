package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/db"
)

func BoardRoutes(r *chi.Mux, database *db.DB) {
	boardHandler := handler.NewBoardHandler(database)

	r.Get("/boards", boardHandler.Index)
	r.Post("/boards", boardHandler.Store)
	r.Get("/boards/{id}", boardHandler.Store)
	r.Put("/boards/{id}", boardHandler.Store)
	r.Delete("/boards/{id}", boardHandler.Store)
}

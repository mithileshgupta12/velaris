package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/db"
)

func PostRoutes(r *chi.Mux, database *db.DB) {
	postHandler := handler.NewPostHandler(database)

	r.Get("/posts", postHandler.Index)
	r.Post("/posts", postHandler.Store)
}

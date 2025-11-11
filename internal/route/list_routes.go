package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/handler"
	"github.com/mithileshgupta12/velaris/internal/middleware"
)

func ListRoutes(r *chi.Mux, listRepository repository.ListRepository, middlewares middleware.Middlewares) {
	listHandler := handler.NewListHandler(listRepository)

	r.Route("/boards/{boardId}/lists", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/", listHandler.Index)
		r.Post("/", listHandler.Store)
		r.Get("/{id}", listHandler.Show)
		r.Put("/{id}", listHandler.Update)
		r.Delete("/{id}", listHandler.Destroy)
	})
}

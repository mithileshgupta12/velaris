package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

func ListRoutes(r *chi.Mux, boardRepository repository.BoardRepository, middlewares middleware.Middlewares) {
	listHandler := handler.NewListHandler()

	r.Route("/boards/{boardId}/lists", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/", listHandler.Index)
		r.Post("/", listHandler.Store)
		r.Get("/{id}", listHandler.Show)
		r.Put("/{id}", listHandler.Update)
		r.Delete("/{id}", listHandler.Destroy)
	})
}

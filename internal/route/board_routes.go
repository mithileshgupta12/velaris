package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/handler"
	"github.com/mithileshgupta12/velaris/internal/middleware"
)

func BoardRoutes(r *chi.Mux, boardRepository repository.BoardRepository, middlewares middleware.Middlewares) {
	boardHandler := handler.NewBoardHandler(boardRepository)

	r.Route("/boards", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/", boardHandler.Index)
		r.Post("/", boardHandler.Store)
		r.Get("/{id}", boardHandler.Show)
		r.Put("/{id}", boardHandler.Update)
		r.Delete("/{id}", boardHandler.Destroy)
	})

}

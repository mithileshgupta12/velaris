package route

import (
	"example.com/velaris/internal/api/handler"
	"github.com/go-chi/chi/v5"
)

func PostRoutes(r *chi.Mux) {
	postHandler := handler.NewPostHandler()
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.Index)
		r.Post("/", postHandler.Store)
		r.Put("/Put{postId}", postHandler.Update)
		r.Delete("/{postId}", postHandler.Update)
	})
}

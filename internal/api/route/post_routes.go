package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/api/handler"
)

func PostRoutes(r *chi.Mux) {
	postHandler := handler.NewPostHandler()

	r.Get("/posts", postHandler.Index)
}

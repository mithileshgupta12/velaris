package route

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func AuthRoutes(r *chi.Mux, queries repository.Querier, lgr logger.Logger) {
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("login successful"))
		if err != nil {
			lgr.Log(logger.ERROR, fmt.Sprintf("some error occurred: %v", err), nil)
		}
	})
}

package helper

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseID(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		return 0, err
	}

	return int64(id), err
}

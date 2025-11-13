package helper

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseIntURLParam(r *http.Request, urlParam string) (int64, error) {
	param := chi.URLParam(r, urlParam)

	id, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return int64(id), err
}

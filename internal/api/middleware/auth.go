package middleware

import (
	"context"
	"net/http"
)

type ctxUserKey string

const CtxUserKey ctxUserKey = "ctxUser"

type CtxUser struct {
	ID    int
	Name  string
	Email string
}

func AuthMiddleware(next http.Handler) http.Handler {
	ctxUser := CtxUser{
		ID:    1,
		Name:  "Bob",
		Email: "bob@aol.com",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CtxUserKey, ctxUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

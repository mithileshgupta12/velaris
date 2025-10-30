package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/mithileshgupta12/velaris/internal/helper"
)

type ctxUserKey string

const CtxUserKey ctxUserKey = "ctxUser"

type CtxUser struct {
	ID    int
	Name  string
	Email string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("auth_session")
		if err != nil {
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		log.Println(sessionCookie)

		ctxUser := CtxUser{
			ID:    1,
			Name:  "Bob",
			Email: "bob@aol.com",
		}
		ctx := context.WithValue(r.Context(), CtxUserKey, ctxUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

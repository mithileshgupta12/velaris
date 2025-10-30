package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/mithileshgupta12/velaris/internal/helper"
)

type ctxUserKey string

const CtxUserKey ctxUserKey = "ctxUser"

type CtxUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (m *middlewares) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("auth_session")
		if err != nil {
			helper.ErrorJsonResponse(w, http.StatusBadRequest, "Unauthenticated")
			return
		}

		sessionData, err := m.sessionStore.Get(r.Context(), sessionCookie.Value)
		if err != nil {
			helper.ErrorJsonResponse(w, http.StatusBadRequest, "Unauthenticated")
			return
		}

		userId, err := strconv.Atoi(sessionData)
		if err != nil {
			helper.ErrorJsonResponse(w, http.StatusBadRequest, "Unauthenticated")
			return
		}

		user, err := m.queries.GetUserById(r.Context(), int64(userId))
		if err != nil {
			helper.ErrorJsonResponse(w, http.StatusBadRequest, "Unauthenticated")
			return
		}

		ctxUser := CtxUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}

		ctx := context.WithValue(r.Context(), CtxUserKey, ctxUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

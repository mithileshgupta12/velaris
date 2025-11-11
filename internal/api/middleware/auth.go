package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/mithileshgupta12/velaris/internal/helper"
)

type ctxUserKey string

const CtxUserKey ctxUserKey = "ctxUser"

type CtxUser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const AuthCookieName = "auth_session"

func (m *middlewares) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		isSecure := r.TLS != nil

		sessionData, err := m.sessionStore.Get(r.Context(), sessionCookie.Value)
		if err != nil {
			helper.SetCookie(w, AuthCookieName, "", -1, isSecure)
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		userId, err := strconv.Atoi(sessionData)
		if err != nil {
			helper.SetCookie(w, AuthCookieName, "", -1, isSecure)
			if err := m.sessionStore.Del(r.Context(), sessionCookie.Value); err != nil {
				slog.Error("failed to delete entry from session store", "err", err)
			}
			slog.Error("failed to convert userId to int", "err", err)
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		user, err := m.repositories.GetUserById(int64(userId))
		if err != nil {
			helper.SetCookie(w, AuthCookieName, "", -1, isSecure)
			if err := m.sessionStore.Del(r.Context(), sessionCookie.Value); err != nil {
				slog.Error("failed to delete entry from session store", "err", err)
			}
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		ctxUser := CtxUser{
			ID:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		ctx := context.WithValue(r.Context(), CtxUserKey, ctxUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
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
				m.lgr.Log(
					logger.ERROR,
					fmt.Sprintf("failed to delete entry from session store: %v", err),
					nil,
				)
			}
			m.lgr.Log(logger.ERROR, fmt.Sprintf("failed to convert userId to int: %v", err), nil)
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		user, err := m.queries.GetUserById(r.Context(), int64(userId))
		if err != nil {
			helper.SetCookie(w, AuthCookieName, "", -1, isSecure)
			if err := m.sessionStore.Del(r.Context(), sessionCookie.Value); err != nil {
				m.lgr.Log(
					logger.ERROR,
					fmt.Sprintf("failed to delete entry from session store: %v", err),
					nil,
				)
			}
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		ctxUser := CtxUser{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		ctx := context.WithValue(r.Context(), CtxUserKey, ctxUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

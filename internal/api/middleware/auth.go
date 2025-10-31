package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

type ctxUserKey string

const CtxUserKey ctxUserKey = "ctxUser"

type CtxUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

const AuthCookieName = "auth_session"

func (m *middlewares) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "Unauthenticated")
			return
		}

		isSecure := r.TLS != nil

		sessionData, err := m.sessionStore.Get(r.Context(), sessionCookie.Value)
		if err != nil {
			helper.SetCookie(w, AuthCookieName, "", -1, isSecure)
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "Unauthenticated")
			return
		}

		userId, err := strconv.Atoi(sessionData)
		if err != nil {
			helper.SetCookie(w, AuthCookieName, "", -1, isSecure)
			if err := m.sessionStore.Del(r.Context(), sessionCookie.Value); err != nil {
				m.lgr.Log(logger.ERROR, fmt.Sprintf("failed to delete entry from session store: %v", err), []*logger.Field{
					{Key: "session_id", Value: sessionCookie.Value},
					{Key: "user_id", Value: sessionData},
				})
			}
			m.lgr.Log(logger.ERROR, fmt.Sprintf("failed to convert userId to int: %v", err), []*logger.Field{
				{Key: "user_id", Value: sessionData},
			})
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "Unauthenticated")
			return
		}

		user, err := m.queries.GetUserById(r.Context(), int64(userId))
		if err != nil {
			helper.SetCookie(w, AuthCookieName, "", -1, isSecure)
			if err := m.sessionStore.Del(r.Context(), sessionCookie.Value); err != nil {
				m.lgr.Log(logger.ERROR, fmt.Sprintf("failed to delete entry from session store: %v", err), []*logger.Field{
					{Key: "session_id", Value: sessionCookie.Value},
					{Key: "user_id", Value: sessionData},
				})
			}
			helper.ErrorJsonResponse(w, http.StatusUnauthorized, "Unauthenticated")
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

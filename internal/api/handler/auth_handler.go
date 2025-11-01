package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

type RegisterUserRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	queries      repository.Querier
	sessionStore cache.SessionStore
	lgr          logger.Logger
}

func NewAuthHandler(queries repository.Querier, sessionStore cache.SessionStore, lgr logger.Logger) *AuthHandler {
	return &AuthHandler{queries, sessionStore, lgr}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&registerUserRequest); err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to decode request: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	registerUserRequest.Name = strings.TrimSpace(registerUserRequest.Name)
	registerUserRequest.Email = strings.ToLower(strings.TrimSpace(registerUserRequest.Email))

	if registerUserRequest.Name == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name is a required field")
		return
	}

	if len(registerUserRequest.Name) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name must not be more than 255 characters long")
		return
	}

	if registerUserRequest.Email == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "email is a required field")
		return
	}

	if len(registerUserRequest.Email) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "email must not be more than 255 characters long")
		return
	}

	if registerUserRequest.Password == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password is a required field")
		return
	}

	if len(registerUserRequest.Password) < 8 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password must be at least 8 characters long")
		return
	}

	if len(registerUserRequest.Password) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password must not be more than 255 characters long")
		return
	}

	if registerUserRequest.PasswordConfirmation == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password_confirmation is a required field")
		return
	}

	if len(registerUserRequest.PasswordConfirmation) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password_confirmation must not be more than 255 characters long")
		return
	}

	if _, err := mail.ParseAddress(registerUserRequest.Email); err != nil {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "email must be a valid email")
		return
	}

	if registerUserRequest.Password != registerUserRequest.PasswordConfirmation {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password and password_confirmation do not match")
		return
	}

	hashedPassword, err := helper.HashPassword(registerUserRequest.Password)
	if err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to hash password: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	user, err := ah.queries.CreateUser(r.Context(), repository.CreateUserParams{
		Name:     registerUserRequest.Name,
		Email:    registerUserRequest.Email,
		Password: hashedPassword,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				helper.ErrorJsonResponse(w, http.StatusConflict, "email is already taken")
				return
			}
		}
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to register user: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusCreated, user)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginUserRequest LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to decode request: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	loginUserRequest.Email = strings.ToLower(strings.TrimSpace(loginUserRequest.Email))

	if loginUserRequest.Email == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "email is a required field")
		return
	}

	if len(loginUserRequest.Email) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "email must not be more than 255 characters long")
		return
	}

	if loginUserRequest.Password == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password is a required field")
		return
	}

	if len(loginUserRequest.Password) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password must not be more than 255 characters long")
		return
	}

	user, err := ah.queries.GetUserByEmail(r.Context(), loginUserRequest.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			helper.ErrorJsonResponse(w, http.StatusBadRequest, "username or password is invalid")
			return
		}

		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to get user by email: %v", err), []*logger.Field{
			{Key: "user_email", Value: loginUserRequest.Email},
		})
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	ok, err := helper.VerifyPassword(loginUserRequest.Password, user.Password)
	if err != nil || !ok {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "username or password is invalid")
		return
	}

	sessionID := make([]byte, 32)
	_, err = rand.Read(sessionID)
	if err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to create session ID: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	b64SessionID := base64.RawURLEncoding.EncodeToString(sessionID)

	if err := ah.sessionStore.Set(r.Context(), b64SessionID, user.ID, time.Duration(time.Hour*24)); err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to set value in session store: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	isSecure := r.TLS != nil
	helper.SetCookie(w, middleware.AuthCookieName, b64SessionID, 60*60*24, isSecure)

	helper.JsonResponse(w, http.StatusOK, "Logged in successfully")
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("auth_session")
	if err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to get session cookie: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusUnauthorized, "Unautenticated")
		return
	}

	if err := ah.sessionStore.Del(r.Context(), sessionCookie.Value); err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to delete record from session: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	isSecure := r.TLS != nil
	helper.SetCookie(w, middleware.AuthCookieName, "", -1, isSecure)

	helper.JsonResponse(w, http.StatusOK, "Logged out successfully")
}

func (ah *AuthHandler) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	loggedInUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	helper.JsonResponse(w, http.StatusOK, loggedInUser)
}

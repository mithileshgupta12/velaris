package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
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
	queries repository.Querier
	lgr     logger.Logger
}

func NewAuthHandler(queries repository.Querier, lgr logger.Logger) *AuthHandler {
	return &AuthHandler{queries, lgr}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&registerUserRequest); err != nil {
		ah.lgr.Log(logger.ERROR, fmt.Sprintf("failed to decode request: %v", err), nil)
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	registerUserRequest.Name = strings.TrimSpace(registerUserRequest.Name)
	registerUserRequest.Email = strings.TrimSpace(registerUserRequest.Email)

	if registerUserRequest.Name == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name is a required field")
		return
	}

	if registerUserRequest.Email == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "email is a required field")
		return
	}

	if registerUserRequest.Password == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password is a required field")
		return
	}

	if registerUserRequest.PasswordConfirmation == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "password_confirmation is a required field")
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
			if pgErr.Code == "23505" {
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
	//
}

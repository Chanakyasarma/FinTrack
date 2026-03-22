package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fintrack/internal/service"
	"fintrack/pkg/respond"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	errs := make(map[string]string)
	if input.Email == "" {
		errs["email"] = "email is required"
	}
	if len(input.Password) < 8 {
		errs["password"] = "password must be at least 8 characters"
	}
	if input.FullName == "" {
		errs["full_name"] = "full name is required"
	}
	if len(errs) > 0 {
		respond.ValidationError(w, errs)
		return
	}

	result, err := h.authService.Register(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			respond.Error(w, http.StatusConflict, err.Error())
			return
		}
		respond.Error(w, http.StatusInternalServerError, "registration failed")
		return
	}

	respond.Success(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Email == "" || input.Password == "" {
		respond.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	result, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respond.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		respond.Error(w, http.StatusInternalServerError, "login failed")
		return
	}

	respond.Success(w, http.StatusOK, result)
}

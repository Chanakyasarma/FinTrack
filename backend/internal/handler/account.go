package handler

import (
	"encoding/json"
	"net/http"

	"fintrack/internal/service"
	"fintrack/pkg/middleware"
	"fintrack/pkg/respond"

	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var input service.CreateAccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	errs := make(map[string]string)
	if input.Name == "" {
		errs["name"] = "name is required"
	}
	if input.Type == "" {
		errs["type"] = "type is required (checking, savings, credit)"
	}
	if len(errs) > 0 {
		respond.ValidationError(w, errs)
		return
	}

	account, err := h.accountService.Create(r.Context(), userID, input)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	respond.Success(w, http.StatusCreated, account)
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	accounts, err := h.accountService.List(r.Context(), userID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to list accounts")
		return
	}

	respond.Success(w, http.StatusOK, accounts)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	accountID := chi.URLParam(r, "id")

	account, err := h.accountService.Get(r.Context(), userID, accountID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, "account not found")
		return
	}

	respond.Success(w, http.StatusOK, account)
}

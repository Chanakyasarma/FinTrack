package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fintrack/internal/service"
	"fintrack/pkg/middleware"
	"fintrack/pkg/respond"

	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	txService *service.TransactionService
}

func NewTransactionHandler(txService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{txService: txService}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var input service.CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	errs := make(map[string]string)
	if input.AccountID == "" {
		errs["account_id"] = "account_id is required"
	}
	if input.Amount <= 0 {
		errs["amount"] = "amount must be greater than 0"
	}
	if input.Type == "" {
		errs["type"] = "type is required (credit or debit)"
	}
	if len(errs) > 0 {
		respond.ValidationError(w, errs)
		return
	}

	if input.Category == "" {
		input.Category = "other"
	}

	tx, err := h.txService.Create(r.Context(), userID, input)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}

	respond.Success(w, http.StatusCreated, tx)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	txs, err := h.txService.List(r.Context(), userID, limit, offset)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to list transactions")
		return
	}

	respond.Success(w, http.StatusOK, txs)
}

func (h *TransactionHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	summary, err := h.txService.Summary(r.Context(), userID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get summary")
		return
	}

	respond.Success(w, http.StatusOK, summary)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	txID := chi.URLParam(r, "id")

	var input service.UpdateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	errs := make(map[string]string)
	if input.Amount <= 0 {
		errs["amount"] = "amount must be greater than 0"
	}
	if input.Type == "" {
		errs["type"] = "type is required (credit or debit)"
	}
	if len(errs) > 0 {
		respond.ValidationError(w, errs)
		return
	}

	if input.Category == "" {
		input.Category = "other"
	}

	tx, err := h.txService.Update(r.Context(), userID, txID, input)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.Success(w, http.StatusOK, tx)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	txID := chi.URLParam(r, "id")

	if err := h.txService.Delete(r.Context(), userID, txID); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.Success(w, http.StatusOK, map[string]string{"message": "transaction deleted"})
}

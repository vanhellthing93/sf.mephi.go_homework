package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.ParseUint(vars["account_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	var request struct {
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"` // "income" или "expense"
		Category    string  `json:"category"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что счет принадлежит пользователю
	if !h.service.AccountBelongsToUser(uint(accountID), userID) {
		http.Error(w, "Account does not belong to user", http.StatusForbidden)
		return
	}

	transaction, err := h.service.CreateTransaction(uint(accountID), request.Amount, request.Type, request.Category, request.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.ParseUint(vars["account_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Проверяем, что счет принадлежит пользователю
	if !h.service.AccountBelongsToUser(uint(accountID), userID) {
		http.Error(w, "Account does not belong to user", http.StatusForbidden)
		return
	}

	transactions, err := h.service.GetTransactionsByAccountID(uint(accountID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, err := strconv.ParseUint(vars["transaction_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.GetTransactionByID(uint(transactionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Проверяем, что операция принадлежит пользователю
	if !h.service.AccountBelongsToUser(transaction.AccountID, userID) {
		http.Error(w, "Transaction does not belong to user", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, err := strconv.ParseUint(vars["transaction_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"` // "income" или "expense"
		Category    string  `json:"category"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем текущую операцию
	transaction, err := h.service.GetTransactionByID(uint(transactionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Проверяем, что операция принадлежит пользователю
	if !h.service.AccountBelongsToUser(transaction.AccountID, userID) {
		http.Error(w, "Transaction does not belong to user", http.StatusForbidden)
		return
	}

	// Обновляем операцию
	transaction.Amount = request.Amount
	transaction.Type = request.Type
	transaction.Category = request.Category
	transaction.Description = request.Description

	if err := h.service.UpdateTransaction(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, err := strconv.ParseUint(vars["transaction_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	// Получаем текущую операцию
	transaction, err := h.service.GetTransactionByID(uint(transactionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Проверяем, что операция принадлежит пользователю
	if !h.service.AccountBelongsToUser(transaction.AccountID, userID) {
		http.Error(w, "Transaction does not belong to user", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteTransaction(uint(transactionID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

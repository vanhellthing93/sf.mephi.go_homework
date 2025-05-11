package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста как uint
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Currency string `json:"currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := h.service.CreateAccount(userID, request.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	accounts, err := h.service.GetAccountsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(accounts)
}
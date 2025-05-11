// internal/handlers/credit_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

type CreditHandler struct {
	service *services.CreditService
}

func NewCreditHandler(service *services.CreditService) *CreditHandler {
	return &CreditHandler{service: service}
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	var request struct {
		Amount float64 `json:"amount"`
		Term   int     `json:"term"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	credit, err := h.service.CreateCredit(userID, request.Amount, request.Term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetUserCredits(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	credits, err := h.service.GetCreditsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credits)
}

func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditID, err := strconv.ParseUint(vars["credit_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Проверяем, что кредит принадлежит пользователю
	if !h.service.CreditBelongsToUser(uint(creditID), userID) {
		http.Error(w, "Credit does not belong to user", http.StatusForbidden)
		return
	}

	schedule, err := h.service.GetPaymentSchedule(uint(creditID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(schedule)
}

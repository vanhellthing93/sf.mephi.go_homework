package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditID, err := strconv.ParseUint(vars["credit_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	var request struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что кредит принадлежит пользователю
	if !h.service.CreditBelongsToUser(uint(creditID), userID) {
		http.Error(w, "Credit does not belong to user", http.StatusForbidden)
		return
	}

	payment, err := h.service.CreatePayment(uint(creditID), request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) GetCreditPayments(w http.ResponseWriter, r *http.Request) {
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

	payments, err := h.service.GetPaymentsByCreditID(uint(creditID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payments)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID, err := strconv.ParseUint(vars["payment_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	payment, err := h.service.GetPaymentByID(uint(paymentID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем, что кредит принадлежит пользователю
	if !h.service.CreditBelongsToUser(payment.CreditID, userID) {
		http.Error(w, "Credit does not belong to user", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

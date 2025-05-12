package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/utils"
)

type CardHandler struct {
	service *services.CardService
}

func NewCardHandler(service *services.CardService) *CardHandler {
	return &CardHandler{service: service}
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.ParseUint(vars["account_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Проверяем, что аккаунт принадлежит пользователю
	if !h.service.AccountBelongsToUser(uint(accountID), userID) {
		http.Error(w, "Account does not belong to user", http.StatusForbidden)
		return
	}

	card, err := h.service.CreateCard(uint(accountID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) GetAccountCards(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    accountID, err := strconv.ParseUint(vars["account_id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid account ID", http.StatusBadRequest)
        return
    }

    // Получаем userID из контекста
    userID := r.Context().Value("userID").(uint)

    // Проверяем, что аккаунт принадлежит пользователю
    if !h.service.AccountBelongsToUser(uint(accountID), userID) {
        http.Error(w, "Account does not belong to user", http.StatusForbidden)
        return
    }

    cards, err := h.service.GetCardsByAccountID(uint(accountID))
    if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Error getting cards")        
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(cards)
}

func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID, err := strconv.ParseUint(vars["card_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Получаем карту
	card, err := h.service.GetCardByID(uint(cardID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем, что карта принадлежит пользователю
	if !h.service.CardBelongsToUser(uint(cardID), userID) {
		http.Error(w, "Card does not belong to user", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID, err := strconv.ParseUint(vars["card_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uint)

	// Проверяем, что карта принадлежит пользователю
	if !h.service.CardBelongsToUser(uint(cardID), userID) {
		http.Error(w, "Card does not belong to user", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteCard(uint(cardID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

type TransferHandler struct {
    service *services.TransferService
}

func NewTransferHandler(service *services.TransferService) *TransferHandler {
    return &TransferHandler{service: service}
}

func (h *TransferHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fromAccountID, err := strconv.ParseUint(vars["from_account_id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid from account ID", http.StatusBadRequest)
        return
    }

    // Получаем userID из контекста
    userID := r.Context().Value("userID").(uint)

    log.Printf("Creating transfer from account %d by user %d", fromAccountID, userID)

    var request struct {
        ToAccount   uint    `json:"to_account"`
        Amount      float64 `json:"amount"`
        Description string  `json:"description"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Transfer request: to_account=%d, amount=%.2f, description=%s", request.ToAccount, request.Amount, request.Description)

    // Проверяем, что счет отправителя принадлежит пользователю
    if !h.service.AccountBelongsToUser(uint(fromAccountID), userID) {
        log.Printf("Account %d does not belong to user %d", fromAccountID, userID)
        http.Error(w, "From account does not belong to user", http.StatusForbidden)
        return
    }

    transfer, err := h.service.CreateTransfer(uint(fromAccountID), request.ToAccount, request.Amount, request.Description)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(transfer)
}

func (h *TransferHandler) GetAccountTransfers(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    accountID, err := strconv.ParseUint(vars["account_id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid account ID", http.StatusBadRequest)
        return
    }

    userID, ok := r.Context().Value("userID").(uint)
    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    if !h.service.AccountBelongsToUser(uint(accountID), userID) {
        http.Error(w, "Account does not belong to user", http.StatusForbidden)
        return
    }

    transfers, err := h.service.GetTransfersByAccountID(uint(accountID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := json.NewEncoder(w).Encode(transfers); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func (h *TransferHandler) GetTransfer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    transferID, err := strconv.ParseUint(vars["transfer_id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid transfer ID", http.StatusBadRequest)
        return
    }

    transfer, err := h.service.GetTransferByID(uint(transferID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if transfer == nil {
        http.Error(w, "Transfer not found", http.StatusNotFound)
        return
    }

    userID, ok := r.Context().Value("userID").(uint)
    if !ok {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

    if !h.service.TransferBelongsToUser(uint(transferID), userID) {
        http.Error(w, "Transfer does not belong to user", http.StatusForbidden)
        return
    }

    if err := json.NewEncoder(w).Encode(transfer); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}
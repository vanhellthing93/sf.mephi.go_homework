package models

import "time"

type Card struct {
	ID        uint      `json:"id"`
	AccountID uint      `json:"account_id"`
	Number    string    `json:"number"`
	CVV       string    `json:"-"` // Не возвращаем CVV в ответе
	Expiry    string    `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
	HMAC      string    `json:"-"` // Не возвращаем HMAC в ответе
}
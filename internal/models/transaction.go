package models

import "time"

type Transaction struct {
	ID          uint      `json:"id"`
	AccountID    uint      `json:"account_id"`
	Amount       float64   `json:"amount"`
	Type         string    `json:"type"` // "income" или "expense"
	Category      string    `json:"category"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

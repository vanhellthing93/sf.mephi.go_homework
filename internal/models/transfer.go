package models

import "time"

type Transfer struct {
	ID          uint      `json:"id"`
	FromAccount  uint      `json:"from_account"`
	ToAccount    uint      `json:"to_account"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}
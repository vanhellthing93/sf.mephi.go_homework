package models

import "time"

type Account struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}
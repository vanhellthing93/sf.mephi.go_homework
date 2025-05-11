package models

import "time"

type Payment struct {
	ID          uint      `json:"id"`
	CreditID     uint      `json:"credit_id"`
	Amount       float64   `json:"amount"`
	PaymentDate  time.Time `json:"payment_date"`
	Status       string    `json:"status"` // "pending", "completed", "failed"
	CreatedAt    time.Time `json:"created_at"`
}
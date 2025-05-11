package models

import "time"

type Credit struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Amount       float64   `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	Term         int       `json:"term"` // в месяцах
	CreatedAt    time.Time `json:"created_at"`
}

type PaymentSchedule struct {
	ID        uint      `json:"id"`
	CreditID  uint      `json:"credit_id"`
	DueDate   time.Time `json:"due_date"`
	Amount    float64   `json:"amount"`
	IsPaid    bool      `json:"is_paid"`
	CreatedAt time.Time `json:"created_at"`
}
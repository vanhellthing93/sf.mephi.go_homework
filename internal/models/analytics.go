package models

import "time"

type IncomeExpenseStats struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type BalanceForecast struct {
	Date    time.Time `json:"date"`
	Balance float64  `json:"balance"`
}

type CreditLoad struct {
	TotalDebt float64 `json:"total_debt"`
	MonthlyPayment float64 `json:"monthly_payment"`
}
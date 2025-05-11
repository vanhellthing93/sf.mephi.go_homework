package repositories

import (
	"database/sql"
	"fmt"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type PaymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) CreatePayment(payment *models.Payment) error {
	query := `INSERT INTO payments (credit_id, amount, payment_date, status, created_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	return r.DB.QueryRow(query, payment.CreditID, payment.Amount, payment.PaymentDate, payment.Status, payment.CreatedAt).Scan(&payment.ID)
}

func (r *PaymentRepository) GetPaymentsByCreditID(creditID uint) ([]models.Payment, error) {
	var payments []models.Payment
	query := `SELECT id, credit_id, amount, payment_date, status, created_at
	          FROM payments
	          WHERE credit_id=$1
	          ORDER BY payment_date DESC`
	rows, err := r.DB.Query(query, creditID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.ID, &payment.CreditID, &payment.Amount, &payment.PaymentDate, &payment.Status, &payment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %v", err)
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *PaymentRepository) GetPaymentByID(paymentID uint) (*models.Payment, error) {
	var payment models.Payment
	query := `SELECT id, credit_id, amount, payment_date, status, created_at
	          FROM payments
	          WHERE id=$1`
	err := r.DB.QueryRow(query, paymentID).Scan(&payment.ID, &payment.CreditID, &payment.Amount, &payment.PaymentDate, &payment.Status, &payment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdatePaymentStatus(paymentID uint, status string) error {
	query := `UPDATE payments SET status=$1 WHERE id=$2`
	_, err := r.DB.Exec(query, status, paymentID)
	return err
}

func (r *PaymentRepository) GetOverduePayments() ([]models.Payment, error) {
	var payments []models.Payment
	query := `SELECT id, credit_id, amount, payment_date, status, created_at
	          FROM payments
	          WHERE payment_date < NOW() AND status='pending'`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue payments: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.ID, &payment.CreditID, &payment.Amount, &payment.PaymentDate, &payment.Status, &payment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan overdue payment: %v", err)
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

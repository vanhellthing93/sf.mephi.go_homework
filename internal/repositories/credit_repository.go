package repositories

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type CreditRepository struct {
	DB *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{DB: db}
}

func (r *CreditRepository) CreateCredit(credit *models.Credit) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Создание записи о кредите
	query := `INSERT INTO credits (user_id, amount, interest_rate, term, created_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = tx.QueryRow(query, credit.UserID, credit.Amount, credit.InterestRate, credit.Term, credit.CreatedAt).Scan(&credit.ID)
	if err != nil {
		return fmt.Errorf("failed to create credit: %v", err)
	}

	// Создание графика платежей
	schedule := generatePaymentSchedule(credit)
	for _, payment := range schedule {
		query = `INSERT INTO payment_schedules (credit_id, due_date, amount, is_paid, created_at)
		          VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(query, credit.ID, payment.DueDate, payment.Amount, payment.IsPaid, payment.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create payment schedule: %v", err)
		}
	}

	// Создание платежей со статусом "pending"
	for _, payment := range schedule {
		query = `INSERT INTO payments (credit_id, amount, payment_date, status, created_at)
		          VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(query, credit.ID, payment.Amount, payment.DueDate, "pending", time.Now())
		if err != nil {
			return fmt.Errorf("failed to create pending payment: %v", err)
		}
	}

	return tx.Commit()
}


func (r *CreditRepository) GetCreditsByUserID(userID uint) ([]models.Credit, error) {
	var credits []models.Credit
	query := `SELECT id, user_id, amount, interest_rate, term, created_at
	          FROM credits
	          WHERE user_id=$1
	          ORDER BY created_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credits: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var credit models.Credit
		if err := rows.Scan(&credit.ID, &credit.UserID, &credit.Amount, &credit.InterestRate, &credit.Term, &credit.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan credit: %v", err)
		}
		credits = append(credits, credit)
	}
	return credits, nil
}

func (r *CreditRepository) GetCreditByID(creditID uint) (*models.Credit, error) {
	var credit models.Credit
	query := `SELECT id, user_id, amount, interest_rate, term, created_at
	          FROM credits
	          WHERE id=$1`
	err := r.DB.QueryRow(query, creditID).Scan(&credit.ID, &credit.UserID, &credit.Amount, &credit.InterestRate, &credit.Term, &credit.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit: %v", err)
	}
	return &credit, nil
}

func (r *CreditRepository) GetPaymentSchedule(creditID uint) ([]models.PaymentSchedule, error) {
	var schedule []models.PaymentSchedule
	query := `SELECT id, credit_id, due_date, amount, is_paid, created_at
	          FROM payment_schedules
	          WHERE credit_id=$1
	          ORDER BY due_date ASC`
	rows, err := r.DB.Query(query, creditID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment schedule: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.PaymentSchedule
		if err := rows.Scan(&payment.ID, &payment.CreditID, &payment.DueDate, &payment.Amount, &payment.IsPaid, &payment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan payment schedule: %v", err)
		}
		schedule = append(schedule, payment)
	}
	return schedule, nil
}

func generatePaymentSchedule(credit *models.Credit) []models.PaymentSchedule {
	var schedule []models.PaymentSchedule
	monthlyPayment := calculateMonthlyPayment(credit.Amount, credit.InterestRate, credit.Term)
	for i := 0; i < credit.Term; i++ {
		dueDate := time.Now().AddDate(0, i+1, 0)
		schedule = append(schedule, models.PaymentSchedule{
			CreditID:  credit.ID,
			DueDate:   dueDate,
			Amount:    monthlyPayment,
			IsPaid:    false,
			CreatedAt: time.Now(),
		})
	}
	return schedule
}

func calculateMonthlyPayment(amount, interestRate float64, term int) float64 {
	monthlyRate := interestRate / 12 / 100
	return amount * monthlyRate * math.Pow(1+monthlyRate, float64(term)) / (math.Pow(1+monthlyRate, float64(term)) - 1)
}

func (r *CreditRepository) UpdatePaymentScheduleStatus(paymentID uint, isPaid bool) error {
	query := `UPDATE payment_schedules SET is_paid=$1 WHERE id=$2`
	_, err := r.DB.Exec(query, isPaid, paymentID)
	return err
}
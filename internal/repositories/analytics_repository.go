package repositories

import (
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type AnalyticsRepository struct {
	DB *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{DB: db}
}

func (r *AnalyticsRepository) GetIncomeExpenseStats(userID uint, startDate, endDate time.Time) (*models.IncomeExpenseStats, error) {
	var stats models.IncomeExpenseStats

	// Получаем доходы
	incomeQuery := `SELECT COALESCE(SUM(amount), 0)
	               FROM transfers
	               WHERE to_account IN (SELECT id FROM accounts WHERE user_id=$1)
	               AND created_at BETWEEN $2 AND $3`
	err := r.DB.QueryRow(incomeQuery, userID, startDate, endDate).Scan(&stats.Income)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %v", err)
	}

	// Получаем расходы
	expenseQuery := `SELECT COALESCE(SUM(amount), 0)
	               FROM transfers
	               WHERE from_account IN (SELECT id FROM accounts WHERE user_id=$1)
	               AND created_at BETWEEN $2 AND $3`
	err = r.DB.QueryRow(expenseQuery, userID, startDate, endDate).Scan(&stats.Expense)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense: %v", err)
	}

	return &stats, nil
}

func (r *AnalyticsRepository) GetBalanceForecast(userID uint, days int) ([]models.BalanceForecast, error) {
	var forecast []models.BalanceForecast

	// Получаем счета пользователя
	accounts, err := r.getUserAccounts(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %v", err)
	}
	accountIDs := make([]uint, len(accounts))
	for i, account := range accounts {
		accountIDs[i] = account.ID
	}

	// Получаем переводы пользователя
	transfers, err := r.GetUserTransfers(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transfers: %v", err)
	}

	// Получаем кредиты пользователя
	credits, err := r.getUserCredits(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %v", err)
	}

	// Получаем графики платежей по кредитам
	paymentSchedules, err := r.getPaymentSchedules(credits)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment schedules: %v", err)
	}

	// Вычисляем прогноз баланса
	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, i)
		balance := 0.0

		// Добавляем балансы всех счетов
		for _, account := range accounts {
			balance += account.Balance
		}

		// Учитываем переводы


		for _, transfer := range transfers {
			if transfer.CreatedAt.After(date) {
				continue
			}
			if slices.Contains(accountIDs, transfer.FromAccount) {
				balance -= transfer.Amount
			} else {
				balance += transfer.Amount
			}
		}

		// Учитываем платежи по кредитам
		for _, schedule := range paymentSchedules {
			if schedule.DueDate.After(date) {
				continue
			}
			balance -= schedule.Amount
		}

		forecast = append(forecast, models.BalanceForecast{
			Date:    date,
			Balance: balance,
		})
	}

	return forecast, nil
}

func (r *AnalyticsRepository) GetCreditLoad(userID uint) (*models.CreditLoad, error) {
	var load models.CreditLoad

	// Получаем кредиты пользователя
	credits, err := r.getUserCredits(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %v", err)
	}

	// Вычисляем общую задолженность
	for _, credit := range credits {
		load.TotalDebt += credit.Amount
	}

	// Вычисляем средний ежемесячный платеж
	if len(credits) > 0 {
		load.MonthlyPayment = load.TotalDebt / float64(len(credits))
	}

	return &load, nil
}

func (r *AnalyticsRepository) getUserAccounts(userID uint) ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT id, user_id, balance, currency, created_at
	          FROM accounts
	          WHERE user_id=$1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan account: %v", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *AnalyticsRepository) GetUserTransfers(userID uint) ([]models.Transfer, error) {
	var transfers []models.Transfer
	query := `SELECT id, from_account, to_account, amount, description, created_at
	          FROM transfers
	          WHERE from_account IN (SELECT id FROM accounts WHERE user_id=$1)
	          OR to_account IN (SELECT id FROM accounts WHERE user_id=$1)`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transfers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var transfer models.Transfer
		if err := rows.Scan(&transfer.ID, &transfer.FromAccount, &transfer.ToAccount, &transfer.Amount, &transfer.Description, &transfer.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan transfer: %v", err)
		}
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

func (r *AnalyticsRepository) getUserCredits(userID uint) ([]models.Credit, error) {
	var credits []models.Credit
	query := `SELECT id, user_id, amount, interest_rate, term, created_at
	          FROM credits
	          WHERE user_id=$1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %v", err)
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

func (r *AnalyticsRepository) getPaymentSchedules(credits []models.Credit) ([]models.PaymentSchedule, error) {
	var schedules []models.PaymentSchedule

	for _, credit := range credits {
		var schedule []models.PaymentSchedule
		query := `SELECT id, credit_id, due_date, amount, is_paid, created_at
		          FROM payment_schedules
		          WHERE credit_id=$1`
		rows, err := r.DB.Query(query, credit.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get payment schedules: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var payment models.PaymentSchedule
			if err := rows.Scan(&payment.ID, &payment.CreditID, &payment.DueDate, &payment.Amount, &payment.IsPaid, &payment.CreatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan payment schedule: %v", err)
			}
			schedule = append(schedule, payment)
		}
		schedules = append(schedules, schedule...)
	}

	return schedules, nil
}

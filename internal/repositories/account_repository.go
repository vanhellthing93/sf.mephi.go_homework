package repositories

import (
	"database/sql"
	"fmt"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) CreateAccount(account *models.Account) error {
	// Используем RETURNING для получения ID созданной записи
	query := `INSERT INTO accounts (user_id, balance, currency, created_at)
	          VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.DB.QueryRow(query, account.UserID, account.Balance, account.Currency, account.CreatedAt).Scan(&account.ID)
	return err
}

func (r *AccountRepository) GetAccountsByUserID(userID uint) ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE user_id=$1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *AccountRepository) GetAccountByID(accountID uint) (*models.Account, error) {
    var account models.Account
    query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE id=$1`
    err := r.DB.QueryRow(query, accountID).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to get account: %v", err)
    }
    return &account, nil
}
// internal/repositories/transaction_repository.go
package repositories

import (
	"database/sql"
	"fmt"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	query := `INSERT INTO transactions (account_id, amount, type, category, description, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	return r.DB.QueryRow(query, transaction.AccountID, transaction.Amount, transaction.Type, transaction.Category, transaction.Description, transaction.CreatedAt).Scan(&transaction.ID)
}

func (r *TransactionRepository) GetTransactionsByAccountID(accountID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := `SELECT id, account_id, amount, type, category, description, created_at
	          FROM transactions
	          WHERE account_id=$1
	          ORDER BY created_at DESC`
	rows, err := r.DB.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Type, &transaction.Category, &transaction.Description, &transaction.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTransactionByID(transactionID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	query := `SELECT id, account_id, amount, type, category, description, created_at
	          FROM transactions
	          WHERE id=$1`
	err := r.DB.QueryRow(query, transactionID).Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Type, &transaction.Category, &transaction.Description, &transaction.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}
	return &transaction, nil
}

func (r *TransactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	query := `UPDATE transactions SET amount=$1, type=$2, category=$3, description=$4 WHERE id=$5`
	_, err := r.DB.Exec(query, transaction.Amount, transaction.Type, transaction.Category, transaction.Description, transaction.ID)
	return err
}

func (r *TransactionRepository) DeleteTransaction(transactionID uint) error {
	query := `DELETE FROM transactions WHERE id=$1`
	_, err := r.DB.Exec(query, transactionID)
	return err
}

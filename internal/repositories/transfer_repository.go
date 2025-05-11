package repositories

import (
	"database/sql"
	"fmt"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type TransferRepository struct {
	DB *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{DB: db}
}

func (r *TransferRepository) CreateTransfer(transfer *models.Transfer) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Создание записи о переводе
	query := `INSERT INTO transfers (from_account, to_account, amount, description, created_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = tx.QueryRow(query, transfer.FromAccount, transfer.ToAccount, transfer.Amount, transfer.Description, transfer.CreatedAt).Scan(&transfer.ID)
	if err != nil {
		return fmt.Errorf("failed to create transfer: %v", err)
	}

	// Обновление баланса отправителя
	query = `UPDATE accounts SET balance = balance - $1 WHERE id = $2`
	_, err = tx.Exec(query, transfer.Amount, transfer.FromAccount)
	if err != nil {
		return fmt.Errorf("failed to update sender balance: %v", err)
	}

	// Обновление баланса получателя
	query = `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err = tx.Exec(query, transfer.Amount, transfer.ToAccount)
	if err != nil {
		return fmt.Errorf("failed to update receiver balance: %v", err)
	}

	return tx.Commit()
}

func (r *TransferRepository) GetTransfersByAccountID(accountID uint) ([]models.Transfer, error) {
	var transfers []models.Transfer
	query := `SELECT id, from_account, to_account, amount, description, created_at
	          FROM transfers
	          WHERE from_account=$1 OR to_account=$1
	          ORDER BY created_at DESC`
	rows, err := r.DB.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers: %v", err)
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

func (r *TransferRepository) GetTransferByID(transferID uint) (*models.Transfer, error) {
	var transfer models.Transfer
	query := `SELECT id, from_account, to_account, amount, description, created_at
	          FROM transfers
	          WHERE id=$1`
	err := r.DB.QueryRow(query, transferID).Scan(&transfer.ID, &transfer.FromAccount, &transfer.ToAccount, &transfer.Amount, &transfer.Description, &transfer.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer: %v", err)
	}
	return &transfer, nil
}
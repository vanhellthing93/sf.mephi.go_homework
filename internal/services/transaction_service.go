package services

import (
	"fmt"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
)

type TransactionService struct {
	repo         *repositories.TransactionRepository
	accountRepo  *repositories.AccountRepository
}

func NewTransactionService(repo *repositories.TransactionRepository, accountRepo *repositories.AccountRepository) *TransactionService {
	return &TransactionService{
		repo:         repo,
		accountRepo:  accountRepo,
	}
}

func (s *TransactionService) CreateTransaction(accountID uint, amount float64, transactionType, category, description string) (*models.Transaction, error) {
	// Проверяем, что счет существует и принадлежит пользователю
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %v", err)
	}

	// Создаем операцию
	transaction := &models.Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Type:        transactionType,
		Category:    category,
		Description: description,
		CreatedAt:  time.Now(),
	}

	// Сохраняем операцию в базе данных
	if err := s.repo.CreateTransaction(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Обновляем баланс счета
	if transactionType == "income" {
		account.Balance += amount
	} else {
		account.Balance -= amount
	}

	if err := s.accountRepo.UpdateAccount(account); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %v", err)
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactionsByAccountID(accountID uint) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByAccountID(accountID)
}

func (s *TransactionService) GetTransactionByID(transactionID uint) (*models.Transaction, error) {
	return s.repo.GetTransactionByID(transactionID)
}

func (s *TransactionService) UpdateTransaction(transaction *models.Transaction) error {
	// Получаем текущую операцию
	currentTransaction, err := s.repo.GetTransactionByID(transaction.ID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	// Получаем счет
	account, err := s.accountRepo.GetAccountByID(transaction.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	// Корректируем баланс счета
	if currentTransaction.Type == "income" {
		account.Balance -= currentTransaction.Amount
	} else {
		account.Balance += currentTransaction.Amount
	}

	if transaction.Type == "income" {
		account.Balance += transaction.Amount
	} else {
		account.Balance -= transaction.Amount
	}

	// Обновляем операцию
	if err := s.repo.UpdateTransaction(transaction); err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}

	// Обновляем баланс счета
	if err := s.accountRepo.UpdateAccount(account); err != nil {
		return fmt.Errorf("failed to update account balance: %v", err)
	}

	return nil
}

func (s *TransactionService) DeleteTransaction(transactionID uint) error {
	// Получаем операцию
	transaction, err := s.repo.GetTransactionByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	// Получаем счет
	account, err := s.accountRepo.GetAccountByID(transaction.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	// Корректируем баланс счета
	if transaction.Type == "income" {
		account.Balance -= transaction.Amount
	} else {
		account.Balance += transaction.Amount
	}

	// Удаляем операцию
	if err := s.repo.DeleteTransaction(transactionID); err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}

	// Обновляем баланс счета
	if err := s.accountRepo.UpdateAccount(account); err != nil {
		return fmt.Errorf("failed to update account balance: %v", err)
	}

	return nil
}

func (s *TransactionService) AccountBelongsToUser(accountID, userID uint) bool {
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return false
	}

	return account.UserID == userID
}
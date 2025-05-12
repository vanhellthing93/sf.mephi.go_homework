package services

import (
	"fmt"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
    "github.com/vanhellthing93/sf.mephi.go_homework/internal/utils"
)

type TransferService struct {
	repo          *repositories.TransferRepository
	accountRepo   *repositories.AccountRepository
}

func NewTransferService(repo *repositories.TransferRepository, accountRepo *repositories.AccountRepository) *TransferService {
	return &TransferService{
		repo:          repo,
		accountRepo:   accountRepo,
	}
}

func (s *TransferService) CreateTransfer(fromAccountID, toAccountID uint, amount float64, description string) (*models.Transfer, error) {
    // Проверяем, что счет отправителя существует и принадлежит пользователю
    fromAccount, err := s.accountRepo.GetAccountByID(fromAccountID)
    if err != nil {
        return nil, fmt.Errorf("from account not found: %v", err)
    }

    // Проверяем, что счет получателя существует
    toAccount, err := s.accountRepo.GetAccountByID(toAccountID)
    if err != nil {
        return nil, fmt.Errorf("to account not found: %v", err)
    }

    // Проверяем баланс отправителя
    if fromAccount.Balance < amount {
        return nil, fmt.Errorf("insufficient funds")
    }

    // Проверяем, что валюты совпадают
    if fromAccount.Currency != toAccount.Currency {
        return nil, fmt.Errorf("currency mismatch")
    }

    // Создаем перевод
    transfer := &models.Transfer{
        FromAccount: fromAccountID,
        ToAccount:   toAccountID,
        Amount:      amount,
        Description: description,
        CreatedAt:   time.Now(),
    }

    // Сохраняем перевод в базе данных
    if err := s.repo.CreateTransfer(transfer); err != nil {
        return nil, fmt.Errorf("failed to create transfer: %v", err)
    }

    return transfer, nil
}

func (s *TransferService) GetTransfersByAccountID(accountID uint) ([]models.Transfer, error) {
	return s.repo.GetTransfersByAccountID(accountID)
}

func (s *TransferService) GetTransferByID(transferID uint) (*models.Transfer, error) {
	return s.repo.GetTransferByID(transferID)
}

func (s *TransferService) AccountBelongsToUser(accountID, userID uint) bool {
    account, err := s.accountRepo.GetAccountByID(accountID)
    if err != nil {
        utils.Log.WithError(err).Warn("Error getting account")
        return false
    }

    return account.UserID == userID
}

func (s *TransferService) TransferBelongsToUser(transferID, userID uint) bool {
	transfer, err := s.repo.GetTransferByID(transferID)
	if err != nil {
		return false
	}

	// Проверяем, что перевод принадлежит пользователю
	// Либо отправитель, либо получатель - это пользователь
	fromAccount, err := s.accountRepo.GetAccountByID(transfer.FromAccount)
	if err != nil {
		return false
	}

	toAccount, err := s.accountRepo.GetAccountByID(transfer.ToAccount)
	if err != nil {
		return false
	}

	return fromAccount.UserID == userID || toAccount.UserID == userID
}

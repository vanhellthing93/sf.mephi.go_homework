package services

import (
	"time"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
)

type AccountService struct {
	repo *repositories.AccountRepository
}

func NewAccountService(repo *repositories.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(userID uint, currency string) (*models.Account, error) {
	account := &models.Account{
		UserID:    userID,
		Balance:   0,
		Currency:  currency,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateAccount(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetAccountsByUserID(userID uint) ([]models.Account, error) {
	return s.repo.GetAccountsByUserID(userID)
}

func (s *AccountService) GetAccountByID(accountID uint) (*models.Account, error) {
    return s.repo.GetAccountByID(accountID)
}
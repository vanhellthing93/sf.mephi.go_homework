// internal/services/credit_service.go
package services

import (
	"fmt"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
)

type CreditService struct {
	repo *repositories.CreditRepository
}

func NewCreditService(repo *repositories.CreditRepository) *CreditService {
	return &CreditService{repo: repo}
}

func (s *CreditService) CreateCredit(userID uint, amount, interestRate float64, term int) (*models.Credit, error) {
	credit := &models.Credit{
		UserID:       userID,
		Amount:       amount,
		InterestRate: interestRate,
		Term:         term,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateCredit(credit); err != nil {
		return nil, fmt.Errorf("failed to create credit: %v", err)
	}

	return credit, nil
}

func (s *CreditService) GetCreditsByUserID(userID uint) ([]models.Credit, error) {
	return s.repo.GetCreditsByUserID(userID)
}

func (s *CreditService) GetCreditByID(creditID uint) (*models.Credit, error) {
	return s.repo.GetCreditByID(creditID)
}

func (s *CreditService) GetPaymentSchedule(creditID uint) ([]models.PaymentSchedule, error) {
	return s.repo.GetPaymentSchedule(creditID)
}

func (s *CreditService) CreditBelongsToUser(creditID, userID uint) bool {
	credit, err := s.repo.GetCreditByID(creditID)
	if err != nil {
		return false
	}

	return credit.UserID == userID
}

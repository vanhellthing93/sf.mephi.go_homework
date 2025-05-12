// internal/services/credit_service.go
package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/utils"
)

type CreditService struct {
	repo         *repositories.CreditRepository
	userRepo 	 *repositories.UserRepository
	cbrService    *CBRService
	smtpService   *SMTPService
}

func NewCreditService(repo *repositories.CreditRepository, userRepo *repositories.UserRepository, cbrService *CBRService, smtpService *SMTPService) *CreditService {
	return &CreditService{
		repo:         repo,
		userRepo: 	  userRepo,
		cbrService:    cbrService,
		smtpService:   smtpService,
	}
}

func (s *CreditService) CreateCredit(userID uint, amount float64, term int) (*models.Credit, error) {
	// Получаем константу увеличения ставки из .env
	incrementStr := os.Getenv("CB_RATE_INCREMENT")
	increment, err := strconv.ParseFloat(incrementStr, 64)
	if err != nil {
		increment = 2.5 // Значение по умолчанию
	}
	
	// Получаем текущую ставку ЦБ РФ
	rate, err := s.cbrService.GetCentralBankRate()
	rate += increment
	if err != nil {
		return nil, fmt.Errorf("failed to get CBR rate: %v", err)
	}

	credit := &models.Credit{
		UserID:       userID,
		Amount:       amount,
		InterestRate: rate,
		Term:         term,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateCredit(credit); err != nil {
		return nil, fmt.Errorf("failed to create credit: %v", err)
	}

	// Получаем пользователя
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Отправляем уведомление
	if err := s.smtpService.SendCreditNotification(user.Email, amount, rate, term); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Failed to send credit notification")
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

func (s *CreditService) GetCreditBalance(creditID uint) (float64, error) {
	credit, err := s.repo.GetCreditByID(creditID)
	if err != nil {
		return 0, fmt.Errorf("failed to get credit: %v", err)
	}

	return credit.Amount, nil
}

func (s *CreditService) GetNextPaymentDate(creditID uint) (time.Time, error) {
	schedule, err := s.repo.GetPaymentSchedule(creditID)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get payment schedule: %v", err)
	}

	for _, payment := range schedule {
		if !payment.IsPaid {
			return payment.DueDate, nil
		}
	}

	return time.Time{}, fmt.Errorf("no pending payments found")
}
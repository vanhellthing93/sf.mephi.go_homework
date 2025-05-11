package services

import (
	"fmt"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
)

type PaymentService struct {
	repo         *repositories.PaymentRepository
	creditRepo   *repositories.CreditRepository
	accountRepo  *repositories.AccountRepository
}

func NewPaymentService(repo *repositories.PaymentRepository, creditRepo *repositories.CreditRepository, accountRepo *repositories.AccountRepository) *PaymentService {
	return &PaymentService{
		repo:         repo,
		creditRepo:   creditRepo,
		accountRepo:  accountRepo,
	}
}

func (s *PaymentService) CreatePayment(creditID uint, amount float64) (*models.Payment, error) {
	// Получаем кредит
	credit, err := s.creditRepo.GetCreditByID(creditID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit: %v", err)
	}

	// Проверяем минимальную сумму платежа
	minPayment := 100.0
	if amount < minPayment {
		return nil, fmt.Errorf("payment amount is too small")
	}

	// Проверяем максимальную сумму платежа
	maxPayment := credit.Amount
	if amount > maxPayment {
		return nil, fmt.Errorf("payment amount exceeds credit balance")
	}

	// Проверяем статус кредита
	if credit.Amount <= 0 {
		return nil, fmt.Errorf("credit is already paid off")
	}

	// Создаем платеж
	payment := &models.Payment{
		CreditID:    creditID,
		Amount:      amount,
		PaymentDate: time.Now(),
		Status:      "completed",
		CreatedAt:   time.Now(),
	}

	// Сохраняем платеж в базе данных
	if err := s.repo.CreatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	// Обновляем статус платежа в графике платежей
	if err := s.updatePaymentSchedule(creditID, amount); err != nil {
		return nil, fmt.Errorf("failed to update payment schedule: %v", err)
	}

	// Обновляем баланс кредита
	if err := s.updateCreditBalance(creditID, amount); err != nil {
		return nil, fmt.Errorf("failed to update credit balance: %v", err)
	}

	return payment, nil
}

func (s *PaymentService) GetPaymentsByCreditID(creditID uint) ([]models.Payment, error) {
	return s.repo.GetPaymentsByCreditID(creditID)
}

func (s *PaymentService) GetPaymentByID(paymentID uint) (*models.Payment, error) {
	return s.repo.GetPaymentByID(paymentID)
}

func (s *PaymentService) updatePaymentSchedule(creditID uint, amount float64) error {
	// Получаем график платежей
	schedule, err := s.creditRepo.GetPaymentSchedule(creditID)
	if err != nil {
		return fmt.Errorf("failed to get payment schedule: %v", err)
	}

	// Находим первый непогашенный платеж
	for _, payment := range schedule {
		if !payment.IsPaid {
			// Обновляем статус платежа
			if err := s.creditRepo.UpdatePaymentScheduleStatus(payment.ID, true); err != nil {
				return fmt.Errorf("failed to update payment schedule status: %v", err)
			}
			break
		}
	}

	return nil
}

func (s *PaymentService) updateCreditBalance(creditID uint, amount float64) error {
	// Получаем кредит
	credit, err := s.creditRepo.GetCreditByID(creditID)
	if err != nil {
		return fmt.Errorf("failed to get credit: %v", err)
	}

	// Обновляем баланс кредита
	newBalance := credit.Amount - amount
	query := `UPDATE credits SET amount=$1 WHERE id=$2`
	_, err = s.creditRepo.DB.Exec(query, newBalance, creditID)
	if err != nil {
		return fmt.Errorf("failed to update credit balance: %v", err)
	}

	return nil
}

func (s *PaymentService) ProcessOverduePayments() error {
	// Получаем просроченные платежи
	payments, err := s.repo.GetOverduePayments()
	if err != nil {
		return fmt.Errorf("failed to get overdue payments: %v", err)
	}

	// Обрабатываем каждый просроченный платеж
	for _, payment := range payments {
		// Начисляем штраф
		if err := s.applyPenalty(payment); err != nil {
			return fmt.Errorf("failed to apply penalty: %v", err)
		}

		// Обновляем статус платежа
		if err := s.repo.UpdatePaymentStatus(payment.ID, "failed"); err != nil {
			return fmt.Errorf("failed to update payment status: %v", err)
		}
	}

	return nil
}

func (s *PaymentService) applyPenalty(payment models.Payment) error {
	// Получаем кредит
	credit, err := s.creditRepo.GetCreditByID(payment.CreditID)
	if err != nil {
		return fmt.Errorf("failed to get credit: %v", err)
	}

	// Вычисляем штраф
	penalty := payment.Amount * 0.1 // 10% от суммы платежа

	// Обновляем баланс кредита
	newBalance := credit.Amount + penalty
	query := `UPDATE credits SET amount=$1 WHERE id=$2`
	_, err = s.creditRepo.DB.Exec(query, newBalance, credit.ID)
	if err != nil {
		return fmt.Errorf("failed to update credit balance: %v", err)
	}

	return nil
}

func (s *PaymentService) CreditBelongsToUser(creditID uint, userID uint) bool {
	credit, err := s.creditRepo.GetCreditByID(creditID)
	if err != nil {
		return false
	}

	return credit.UserID == userID
}
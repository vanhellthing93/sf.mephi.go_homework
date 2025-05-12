package services

import (
		"time"
		"github.com/vanhellthing93/sf.mephi.go_homework/internal/utils"
)

type SchedulerService struct {
	paymentService *PaymentService
}

func NewSchedulerService(paymentService *PaymentService) *SchedulerService {
	return &SchedulerService{
		paymentService: paymentService,
	}
}

func (s *SchedulerService) Start() {
	// Запускаем шедулер
	ticker := time.NewTicker(12 * time.Hour)
	go func() {
		s.ProcessOverduePayments()
		for range ticker.C {
			s.ProcessOverduePayments()
		}
	}()
}

func (s *SchedulerService) ProcessOverduePayments() {
	utils.Log.Info("Processing overdue payments...")

	// Обрабатываем просроченные платежи
	if err := s.paymentService.ProcessOverduePayments(); err != nil {
		utils.Log.WithError(err).Warn("Error processing overdue payments")
	}

	utils.Log.Info("Finished processing overdue payments")
}
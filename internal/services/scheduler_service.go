package services

import (
	"log"
	"time"
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
	log.Println("Processing overdue payments...")

	// Обрабатываем просроченные платежи
	if err := s.paymentService.ProcessOverduePayments(); err != nil {
		log.Printf("Error processing overdue payments: %v", err)
	}

	log.Println("Finished processing overdue payments")
}
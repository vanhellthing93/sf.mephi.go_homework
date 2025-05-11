package services

import (
	"fmt"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
)

type AnalyticsService struct {
	repo *repositories.AnalyticsRepository
}

func NewAnalyticsService(repo *repositories.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetIncomeExpenseStats(userID uint, startDate, endDate time.Time) (*models.IncomeExpenseStats, error) {
	return s.repo.GetIncomeExpenseStats(userID, startDate, endDate)
}

func (s *AnalyticsService) GetBalanceForecast(userID uint, days int) ([]models.BalanceForecast, error) {
	return s.repo.GetBalanceForecast(userID, days)
}

func (s *AnalyticsService) GetCreditLoad(userID uint) (*models.CreditLoad, error) {
	return s.repo.GetCreditLoad(userID)
}

func (s *AnalyticsService) GetMonthlyStats(userID uint, year int) ([]models.IncomeExpenseStats, error) {
	var stats []models.IncomeExpenseStats

	for month := 1; month <= 12; month++ {
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

		monthStats, err := s.repo.GetIncomeExpenseStats(userID, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get monthly stats: %v", err)
		}

		stats = append(stats, *monthStats)
	}

	return stats, nil
}

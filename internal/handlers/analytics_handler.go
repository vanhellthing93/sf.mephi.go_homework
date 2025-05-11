package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetIncomeExpenseStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	// Получаем параметры даты из запроса
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // По умолчанию последний месяц
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now()
	}

	stats, err := h.service.GetIncomeExpenseStats(userID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (h *AnalyticsHandler) GetBalanceForecast(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	daysStr := r.URL.Query().Get("days")
	days := 30 // По умолчанию 30 дней

	if daysStr != "" {
		_, err := strconv.Atoi(daysStr)
		if err != nil {
			http.Error(w, "Invalid days parameter", http.StatusBadRequest)
			return
		}
	}

	forecast, err := h.service.GetBalanceForecast(userID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(forecast)
}

func (h *AnalyticsHandler) GetCreditLoad(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	load, err := h.service.GetCreditLoad(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(load)
}

func (h *AnalyticsHandler) GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	yearStr := r.URL.Query().Get("year")
	year := time.Now().Year() // По умолчанию текущий год

	if yearStr != "" {
		_, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid year parameter", http.StatusBadRequest)
			return
		}
	}

	stats, err := h.service.GetMonthlyStats(userID, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

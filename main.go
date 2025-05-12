package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanhellthing93/sf.mephi.go_homework/config"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/handlers"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/middleware"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	if err := config.InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	cardRepo := repositories.NewCardRepository(db)
	transferRepo := repositories.NewTransferRepository(db)
	creditRepo := repositories.NewCreditRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	analyticsRepo := repositories.NewAnalyticsRepository(db)



	// Инициализация сервисов
	smtpService := services.NewSMTPService()
	cbrService := services.NewCBRService()
	userService := services.NewUserService(userRepo, smtpService)
	accountService := services.NewAccountService(accountRepo)
	cardService := services.NewCardService(cardRepo)
	transferService := services.NewTransferService(transferRepo, accountRepo)
	creditService := services.NewCreditService(creditRepo, userRepo, cbrService, smtpService)
	paymentService := services.NewPaymentService(paymentRepo, creditRepo, accountRepo, userRepo, smtpService)
	analyticsService := services.NewAnalyticsService(analyticsRepo)


	// Инициализация шедулера
	schedulerService := services.NewSchedulerService(paymentService)
	schedulerService.Start()

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	cardHandler := handlers.NewCardHandler(cardService)
	transferHandler := handlers.NewTransferHandler(transferService)
	creditHandler := handlers.NewCreditHandler(creditService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)



	r := mux.NewRouter()

	// Публичные маршруты
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Защищенные маршруты
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	// Управление счетами
	authRouter.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	authRouter.HandleFunc("/accounts", accountHandler.GetUserAccounts).Methods("GET")

	// Управление картами
	authRouter.HandleFunc("/accounts/{account_id}/cards", cardHandler.CreateCard).Methods("POST")
	authRouter.HandleFunc("/accounts/{account_id}/cards", cardHandler.GetAccountCards).Methods("GET")
	authRouter.HandleFunc("/cards/{card_id}", cardHandler.GetCard).Methods("GET")
	authRouter.HandleFunc("/cards/{card_id}", cardHandler.DeleteCard).Methods("DELETE")

	// Управление переводами
	authRouter.HandleFunc("/accounts/{from_account_id}/transfers", transferHandler.CreateTransfer).Methods("POST")
	authRouter.HandleFunc("/accounts/{account_id}/transfers", transferHandler.GetAccountTransfers).Methods("GET")
	authRouter.HandleFunc("/transfers/{transfer_id}", transferHandler.GetTransfer).Methods("GET")

	// Управление кредитами
	authRouter.HandleFunc("/credits", creditHandler.CreateCredit).Methods("POST")
	authRouter.HandleFunc("/credits", creditHandler.GetUserCredits).Methods("GET")
	authRouter.HandleFunc("/credits/{credit_id}/schedule", creditHandler.GetPaymentSchedule).Methods("GET")

	// Управление платежами
	authRouter.HandleFunc("/credits/{credit_id}/payments", paymentHandler.CreatePayment).Methods("POST")
	authRouter.HandleFunc("/credits/{credit_id}/payments", paymentHandler.GetCreditPayments).Methods("GET")
	authRouter.HandleFunc("/payments/{payment_id}", paymentHandler.GetPayment).Methods("GET")

	// Аналитика
	authRouter.HandleFunc("/analytics/income-expense", analyticsHandler.GetIncomeExpenseStats).Methods("GET")
	authRouter.HandleFunc("/analytics/balance-forecast", analyticsHandler.GetBalanceForecast).Methods("GET")
	authRouter.HandleFunc("/analytics/credit-load", analyticsHandler.GetCreditLoad).Methods("GET")
	authRouter.HandleFunc("/analytics/monthly-stats", analyticsHandler.GetMonthlyStats).Methods("GET")

	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", r)

}
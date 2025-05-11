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

	// Инициализация сервисов
	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo)
	cardService := services.NewCardService(cardRepo)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	cardHandler := handlers.NewCardHandler(cardService)

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

	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", r)
}
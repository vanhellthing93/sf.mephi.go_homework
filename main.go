package main

import (
	"log"
	"time"
	"github.com/vanhellthing93/sf.mephi.go_homework/config"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/services"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	// Инициализация таблиц
	if err := config.InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	bob := &models.User{
		ID:        2,
		Email:     "bob2@example.com",
		Password:  "securepassword",
		Username:  "bob2",
		CreatedAt: time.Now(),
	}

	err := userService.RegisterUser(bob)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}


}
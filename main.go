package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanhellthing93/sf.mephi.go_homework/config"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/handlers"
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
	userHandler := handlers.NewUserHandler(userService)

	r := mux.NewRouter()

	// Публичные маршруты
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", r)
}
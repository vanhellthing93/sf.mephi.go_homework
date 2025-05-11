package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
)

type UserService struct {
	repo        *repositories.UserRepository
	smtpService *SMTPService
}

func NewUserService(repo *repositories.UserRepository, smtpService *SMTPService) *UserService {
	return &UserService{
		repo:        repo,
		smtpService: smtpService,
	}
}

// Регистрация пользователя
func (s *UserService) RegisterUser(user *models.User) error {
	if err := user.ValidateEmail(); err != nil {
		return err
	}
	if err := user.ValidatePassword(); err != nil {
		return err
	}
	if err := user.HashPassword(); err != nil {
		return err
	}
	user.CreatedAt = time.Now()

	if err := s.repo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	if err := s.smtpService.SendRegistrationNotification(user.Email); err != nil {
		log.Printf("Failed to send registration notification: %v", err)
	}

	return nil
}

// Аутентификация пользователя
func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

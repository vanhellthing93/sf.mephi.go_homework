package models

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// Валидация username
func (u *User) ValidateUsername() error {
    if len(u.Username) < 3 {
        return errors.New("username too short")
    }
    return nil
}



// Валидация email
func (u *User) ValidateEmail() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

// Валидация пароля
func (u *User) ValidatePassword() error {
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

// Хеширование пароля
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Проверка пароля
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

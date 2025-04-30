// internal/repositories/user_repository.go
package repositories

import (
	"database/sql"
	"errors"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Создание пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, password, username, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, user.Email, user.Password, user.Username, user.CreatedAt)
	return err
}

// Получение пользователя по email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password, username, created_at FROM users WHERE email=$1`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

// Получение пользователя по ID
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password, username, created_at FROM users WHERE id=$1`
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

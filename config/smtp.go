package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/utils"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func LoadSMTPConfig() *SMTPConfig {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		utils.Log.Warn("Warning: No .env file found")
	}

	return &SMTPConfig{
		Host:     getEnv("SMTP_HOST", "smtp.example.com"),
		Port:     getPortEnv("SMTP_PORT", 465),
		Username: getEnv("SMTP_USERNAME", "noreply@example.com"),
		Password: getEnv("SMTP_PASSWORD", "password"),
		From:     getEnv("SMTP_FROM", "noreply@example.com"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getPortEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		port, err := strconv.Atoi(value)
		if err != nil {
			utils.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error loading smtp port. Using default values.")
			return defaultValue
		} else {
			return port
		}
	}
	return defaultValue
}

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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
		log.Println("Warning: No .env file found")
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
			log.Printf("Error loading smtp port. Using default values. %e", err)
			return defaultValue
		} else {
			return port
		}
	}
	return defaultValue
}

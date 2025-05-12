package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем параметры подключения из переменных окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("SSL_MODE")

	// Формируем строку подключения
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		sslMode)

	// Открываем соединение с БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to database!")
	return db
}

func InitDB(db *sql.DB) error {
	// Создание таблицы пользователей
	createUsersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		username VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createUsersTableQuery); err != nil {
		return err
	}

	// Создание таблицы счетов
	createAccountsTableQuery := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		balance DECIMAL(15, 2) DEFAULT 0,
		currency VARCHAR(3) DEFAULT 'RUB',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createAccountsTableQuery); err != nil {
		return err
	}

	// Создание таблицы карт
	createCardsTableQuery := `
	CREATE TABLE IF NOT EXISTS cards (
		id SERIAL PRIMARY KEY,
		account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
		number TEXT NOT NULL,
		cvv VARCHAR(60) NOT NULL,
		expiry VARCHAR(5) NOT NULL,
		hmac VARCHAR(64) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createCardsTableQuery); err != nil {
		return err
	}

	// Создание таблицы переводов
	createTransfersTableQuery := `
	CREATE TABLE IF NOT EXISTS transfers (
		id SERIAL PRIMARY KEY,
		from_account INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
		to_account INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
		amount DECIMAL(15, 2) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTransfersTableQuery); err != nil {
		return err
	}

	// Создание таблицы кредитов
	createCreditsTableQuery := `
	CREATE TABLE IF NOT EXISTS credits (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		amount DECIMAL(15, 2) NOT NULL,
		interest_rate DECIMAL(5, 2) NOT NULL,
		term INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createCreditsTableQuery); err != nil {
		return err
	}

	// Создание таблицы графика платежей
	createPaymentSchedulesTableQuery := `
	CREATE TABLE IF NOT EXISTS payment_schedules (
		id SERIAL PRIMARY KEY,
		credit_id INTEGER REFERENCES credits(id) ON DELETE CASCADE,
		due_date TIMESTAMP NOT NULL,
		amount DECIMAL(15, 2) NOT NULL,
		is_paid BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createPaymentSchedulesTableQuery); err != nil {
		return err
	}

	// Создание таблицы платежей
	createPaymentsTableQuery := `
	CREATE TABLE IF NOT EXISTS payments (
		id SERIAL PRIMARY KEY,
		credit_id INTEGER REFERENCES credits(id) ON DELETE CASCADE,
		amount DECIMAL(15, 2) NOT NULL,
		payment_date TIMESTAMP NOT NULL,
		status VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createPaymentsTableQuery); err != nil {
		return err
	}

	// Создание таблицы операций
	createTransactionsTableQuery := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
		amount DECIMAL(15, 2) NOT NULL,
		type VARCHAR(10) NOT NULL,
		category VARCHAR(50),
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTransactionsTableQuery); err != nil {
		return err
	}

	return nil
}
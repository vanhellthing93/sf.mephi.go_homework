package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	connStr := "postgresql://go_hw_admin:gogogogo@192.168.4.173:5433/go_hw?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
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
		number VARCHAR(16) NOT NULL,
		cvv VARCHAR(3) NOT NULL,
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

	return nil
}
package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	log.Println("connected to postgres")
	return db, nil
}

func RunMigrations(db *sqlx.DB) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			full_name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS accounts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			type TEXT NOT NULL CHECK (type IN ('checking','savings','credit')),
			balance NUMERIC(15,2) NOT NULL DEFAULT 0,
			currency TEXT NOT NULL DEFAULT 'USD',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			amount NUMERIC(15,2) NOT NULL,
			type TEXT NOT NULL CHECK (type IN ('credit','debit')),
			category TEXT NOT NULL DEFAULT 'other',
			description TEXT NOT NULL DEFAULT '',
			merchant_name TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("running migration: %w\nSQL: %s", err, m)
		}
	}

	log.Println("migrations applied successfully")
	return nil
}

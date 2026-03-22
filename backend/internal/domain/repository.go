package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id string) (*Account, error)
	ListByUserID(ctx context.Context, userID string) ([]*Account, error)
	UpdateBalance(ctx context.Context, id string, balance float64) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByID(ctx context.Context, id string) (*Transaction, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*Transaction, error)
	ListByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*Transaction, error)
	GetSummary(ctx context.Context, userID string, from, to time.Time) (*Summary, error)
	GetMonthlyTrend(ctx context.Context, userID string, months int) ([]MonthlySummary, error)
}

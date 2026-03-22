package repository

import (
	"context"
	"fmt"

	"fintrack/internal/domain"

	"github.com/jmoiron/sqlx"
)

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) domain.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (user_id, name, type, balance, currency)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	row := r.db.QueryRowxContext(ctx, query,
		account.UserID, account.Name, account.Type, account.Balance, account.Currency)
	return row.Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at FROM accounts WHERE id = $1`
	if err := r.db.GetContext(ctx, &account, query, id); err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return &account, nil
}

func (r *accountRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	var accounts []*domain.Account
	query := `
		SELECT id, user_id, name, type, balance, currency, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at ASC
	`
	if err := r.db.SelectContext(ctx, &accounts, query, userID); err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}
	return accounts, nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, id string, balance float64) error {
	query := `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, balance, id)
	if err != nil {
		return fmt.Errorf("updating balance: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}

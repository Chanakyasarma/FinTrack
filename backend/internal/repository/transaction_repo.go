package repository

import (
	"context"
	"fmt"
	"time"

	"fintrack/internal/domain"

	"github.com/jmoiron/sqlx"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (account_id, user_id, amount, type, category, description, merchant_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	row := r.db.QueryRowxContext(ctx, query,
		tx.AccountID, tx.UserID, tx.Amount, tx.Type, tx.Category, tx.Description, tx.MerchantName)
	return row.Scan(&tx.ID, &tx.CreatedAt)
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var tx domain.Transaction
	query := `SELECT id, account_id, user_id, amount, type, category, description, merchant_name, created_at FROM transactions WHERE id = $1`
	if err := r.db.GetContext(ctx, &tx, query, id); err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &tx, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1, type = $2, category = $3, description = $4, merchant_name = $5
		WHERE id = $6 AND user_id = $7
	`
	result, err := r.db.ExecContext(ctx, query,
		tx.Amount, tx.Type, tx.Category, tx.Description, tx.MerchantName, tx.ID, tx.UserID)
	if err != nil {
		return fmt.Errorf("updating transaction: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("transaction not found or unauthorized")
	}
	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting transaction: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

func (r *transactionRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Transaction, error) {
	var txs []*domain.Transaction
	query := `
		SELECT id, account_id, user_id, amount, type, category, description, merchant_name, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &txs, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("listing transactions: %w", err)
	}
	return txs, nil
}

func (r *transactionRepository) ListByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*domain.Transaction, error) {
	var txs []*domain.Transaction
	query := `
		SELECT id, account_id, user_id, amount, type, category, description, merchant_name, created_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &txs, query, accountID, limit, offset); err != nil {
		return nil, fmt.Errorf("listing transactions: %w", err)
	}
	return txs, nil
}

func (r *transactionRepository) GetSummary(ctx context.Context, userID string, from, to time.Time) (*domain.Summary, error) {
	summary := &domain.Summary{
		CategoryBreakdown: make(map[string]float64),
	}

	// Get total balance across all accounts
	balanceQuery := `SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE user_id = $1`
	if err := r.db.QueryRowContext(ctx, balanceQuery, userID).Scan(&summary.TotalBalance); err != nil {
		return nil, fmt.Errorf("getting total balance: %w", err)
	}

	// Monthly income and expenses
	flowQuery := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'credit' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'debit' THEN amount ELSE 0 END), 0) AS expenses
		FROM transactions
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
	`
	row := r.db.QueryRowContext(ctx, flowQuery, userID, from, to)
	if err := row.Scan(&summary.MonthlyIncome, &summary.MonthlyExpenses); err != nil {
		return nil, fmt.Errorf("getting flow: %w", err)
	}

	// Category breakdown (debits only)
	categoryQuery := `
		SELECT category, COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = $1 AND type = 'debit' AND created_at BETWEEN $2 AND $3
		GROUP BY category
	`
	rows, err := r.db.QueryContext(ctx, categoryQuery, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("getting categories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cat string
		var amount float64
		if err := rows.Scan(&cat, &amount); err != nil {
			continue
		}
		summary.CategoryBreakdown[cat] = amount
	}

	return summary, nil
}

func (r *transactionRepository) GetMonthlyTrend(ctx context.Context, userID string, months int) ([]domain.MonthlySummary, error) {
	query := `
		SELECT
			TO_CHAR(DATE_TRUNC('month', created_at), 'Mon YYYY') AS month,
			COALESCE(SUM(CASE WHEN type = 'credit' THEN amount ELSE 0 END), 0) AS total_credit,
			COALESCE(SUM(CASE WHEN type = 'debit' THEN amount ELSE 0 END), 0) AS total_debit,
			COALESCE(SUM(CASE WHEN type = 'credit' THEN amount ELSE -amount END), 0) AS net_flow
		FROM transactions
		WHERE user_id = $1
			AND created_at >= DATE_TRUNC('month', NOW()) - INTERVAL '1 month' * $2
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at) ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, months)
	if err != nil {
		return nil, fmt.Errorf("getting monthly trend: %w", err)
	}
	defer rows.Close()

	var trend []domain.MonthlySummary
	for rows.Next() {
		var ms domain.MonthlySummary
		if err := rows.Scan(&ms.Month, &ms.TotalCredit, &ms.TotalDebit, &ms.NetFlow); err != nil {
			continue
		}
		trend = append(trend, ms)
	}

	return trend, nil
}

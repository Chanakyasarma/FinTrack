package repository

import (
	"context"
	"fmt"

	"fintrack/internal/domain"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	row := r.db.QueryRowxContext(ctx, query, user.Email, user.PasswordHash, user.FullName)
	return row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, full_name, created_at, updated_at FROM users WHERE email = $1`
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, full_name, created_at, updated_at FROM users WHERE id = $1`
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

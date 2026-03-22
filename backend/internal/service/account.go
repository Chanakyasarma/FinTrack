package service

import (
	"context"
	"fmt"
	"time"

	"fintrack/internal/domain"
)

type AccountService struct {
	accountRepo domain.AccountRepository
	cache       Cache
}

// NewAccountService wires up the service in production.
// Pass a *cache.RedisClient — it satisfies the Cache interface.
func NewAccountService(accountRepo domain.AccountRepository, c Cache) *AccountService {
	return &AccountService{accountRepo: accountRepo, cache: c}
}

// NewAccountServiceWithNilCache is used in unit tests.
func NewAccountServiceWithNilCache(accountRepo domain.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo, cache: &noopCacheImpl{}}
}

type CreateAccountInput struct {
	Name     string             `json:"name"`
	Type     domain.AccountType `json:"type"`
	Currency string             `json:"currency"`
}

func (s *AccountService) Create(ctx context.Context, userID string, input CreateAccountInput) (*domain.Account, error) {
	account := &domain.Account{
		UserID:   userID,
		Name:     input.Name,
		Type:     input.Type,
		Balance:  0,
		Currency: input.Currency,
	}

	if account.Currency == "" {
		account.Currency = "USD"
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	// Invalidate list cache
	s.cache.Delete(ctx, accountListKey(userID))

	return account, nil
}

func (s *AccountService) List(ctx context.Context, userID string) ([]*domain.Account, error) {
	cacheKey := accountListKey(userID)

	var accounts []*domain.Account
	if err := s.cache.Get(ctx, cacheKey, &accounts); err == nil {
		return accounts, nil
	}

	accounts, err := s.accountRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}

	s.cache.Set(ctx, cacheKey, accounts, 5*time.Minute)
	return accounts, nil
}

func (s *AccountService) Get(ctx context.Context, userID, accountID string) (*domain.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}
	if account.UserID != userID {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func accountListKey(userID string) string {
	return fmt.Sprintf("accounts:list:%s", userID)
}

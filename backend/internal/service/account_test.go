package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"fintrack/internal/domain"
	"fintrack/internal/service"
)

// --- Mock AccountRepository ---

type mockAccountRepo struct {
	mu       sync.Mutex
	accounts map[string]*domain.Account
}

func newMockAccountRepo() *mockAccountRepo {
	return &mockAccountRepo{accounts: make(map[string]*domain.Account)}
}

func (m *mockAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.ID = "acct-" + a.Name
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	m.accounts[a.ID] = a
	return nil
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return a, nil
}

func (m *mockAccountRepo) ListByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []*domain.Account
	for _, a := range m.accounts {
		if a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockAccountRepo) UpdateBalance(ctx context.Context, id string, balance float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accounts[id]
	if !ok {
		return errors.New("not found")
	}
	a.Balance = balance
	return nil
}

// --- Tests ---

func TestAccountService_Create(t *testing.T) {
	repo := newMockAccountRepo()
	// Pass nil cache — in real tests, inject a mock cache via interface extraction
	svc := service.NewAccountServiceWithNilCache(repo)

	account, err := svc.Create(context.Background(), "user-1", service.CreateAccountInput{
		Name: "Checking",
		Type: domain.AccountTypeChecking,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account.ID == "" {
		t.Error("expected account to have an ID")
	}
	if account.Currency != "USD" {
		t.Errorf("expected default currency USD, got %s", account.Currency)
	}
	if account.Balance != 0 {
		t.Errorf("expected initial balance 0, got %f", account.Balance)
	}
}

func TestAccountService_Get_WrongUser(t *testing.T) {
	repo := newMockAccountRepo()
	svc := service.NewAccountServiceWithNilCache(repo)

	_, _ = svc.Create(context.Background(), "user-1", service.CreateAccountInput{
		Name: "Savings",
		Type: domain.AccountTypeSavings,
	})

	accounts, _ := repo.ListByUserID(context.Background(), "user-1")
	if len(accounts) == 0 {
		t.Fatal("expected at least one account")
	}

	// Another user tries to access it
	_, err := svc.Get(context.Background(), "user-2", accounts[0].ID)
	if err == nil {
		t.Error("expected error when accessing another user's account")
	}
}

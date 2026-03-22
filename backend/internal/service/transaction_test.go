package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"fintrack/internal/domain"
	"fintrack/internal/service"
	"fintrack/internal/websocket"
)

// --- Mock TransactionRepository ---

type mockTransactionRepo struct {
	mu           sync.Mutex
	transactions map[string]*domain.Transaction
}

func newMockTransactionRepo() *mockTransactionRepo {
	return &mockTransactionRepo{transactions: make(map[string]*domain.Transaction)}
}

func (m *mockTransactionRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	tx.ID = "tx-" + tx.Description
	tx.CreatedAt = time.Now()
	m.transactions[tx.ID] = tx
	return nil
}

func (m *mockTransactionRepo) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tx, ok := m.transactions[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return tx, nil
}

func (m *mockTransactionRepo) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []*domain.Transaction
	for _, tx := range m.transactions {
		if tx.UserID == userID {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockTransactionRepo) ListByAccountID(ctx context.Context, accountID string, limit, offset int) ([]*domain.Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []*domain.Transaction
	for _, tx := range m.transactions {
		if tx.AccountID == accountID {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockTransactionRepo) GetSummary(ctx context.Context, userID string, from, to time.Time) (*domain.Summary, error) {
	return &domain.Summary{
		TotalBalance:      1000,
		MonthlyIncome:     500,
		MonthlyExpenses:   300,
		CategoryBreakdown: map[string]float64{"food": 100, "transport": 200},
	}, nil
}

func (m *mockTransactionRepo) GetMonthlyTrend(ctx context.Context, userID string, months int) ([]domain.MonthlySummary, error) {
	return []domain.MonthlySummary{
		{Month: "Jan 2025", TotalCredit: 3000, TotalDebit: 2000, NetFlow: 1000},
	}, nil
}

// --- Helpers ---

func newTransactionService() (*service.TransactionService, *mockAccountRepo, *mockTransactionRepo) {
	accountRepo := newMockAccountRepo()
	txRepo := newMockTransactionRepo()
	hub := websocket.NewHub()
	go hub.Run()

	svc := service.NewTransactionServiceForTest(txRepo, accountRepo, hub)
	return svc, accountRepo, txRepo
}

// --- Tests ---

func TestTransactionService_Create_Debit(t *testing.T) {
	svc, accountRepo, _ := newTransactionService()

	// Pre-seed an account
	acct := &domain.Account{
		ID:      "acct-1",
		UserID:  "user-1",
		Name:    "Checking",
		Type:    domain.AccountTypeChecking,
		Balance: 1000,
		Currency: "USD",
	}
	accountRepo.mu.Lock()
	accountRepo.accounts[acct.ID] = acct
	accountRepo.mu.Unlock()

	tx, err := svc.Create(context.Background(), "user-1", service.CreateTransactionInput{
		AccountID:    "acct-1",
		Amount:       200,
		Type:         domain.TransactionTypeDebit,
		Category:     domain.CategoryFood,
		Description:  "grocery-run",
		MerchantName: "Whole Foods",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Amount != 200 {
		t.Errorf("expected amount 200, got %f", tx.Amount)
	}

	// Balance should have decreased
	accountRepo.mu.Lock()
	updatedBalance := accountRepo.accounts["acct-1"].Balance
	accountRepo.mu.Unlock()

	if updatedBalance != 800 {
		t.Errorf("expected balance 800, got %f", updatedBalance)
	}
}

func TestTransactionService_Create_Credit(t *testing.T) {
	svc, accountRepo, _ := newTransactionService()

	acct := &domain.Account{
		ID:      "acct-2",
		UserID:  "user-2",
		Name:    "Savings",
		Type:    domain.AccountTypeSavings,
		Balance: 500,
		Currency: "USD",
	}
	accountRepo.mu.Lock()
	accountRepo.accounts[acct.ID] = acct
	accountRepo.mu.Unlock()

	_, err := svc.Create(context.Background(), "user-2", service.CreateTransactionInput{
		AccountID:    "acct-2",
		Amount:       1500,
		Type:         domain.TransactionTypeCredit,
		Category:     domain.CategorySalary,
		Description:  "monthly-salary",
		MerchantName: "ACME Corp",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	accountRepo.mu.Lock()
	updatedBalance := accountRepo.accounts["acct-2"].Balance
	accountRepo.mu.Unlock()

	if updatedBalance != 2000 {
		t.Errorf("expected balance 2000, got %f", updatedBalance)
	}
}

func TestTransactionService_Create_WrongUser(t *testing.T) {
	svc, accountRepo, _ := newTransactionService()

	acct := &domain.Account{
		ID:      "acct-3",
		UserID:  "user-1",
		Balance: 100,
	}
	accountRepo.mu.Lock()
	accountRepo.accounts[acct.ID] = acct
	accountRepo.mu.Unlock()

	_, err := svc.Create(context.Background(), "user-99", service.CreateTransactionInput{
		AccountID: "acct-3",
		Amount:    50,
		Type:      domain.TransactionTypeDebit,
		Category:  domain.CategoryOther,
	})
	if err == nil {
		t.Error("expected error when user does not own the account")
	}
}

func TestTransactionService_Summary(t *testing.T) {
	svc, _, _ := newTransactionService()

	summary, err := svc.Summary(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.MonthlyIncome != 500 {
		t.Errorf("expected income 500, got %f", summary.MonthlyIncome)
	}
	if summary.SavingsRate != 40 {
		t.Errorf("expected savings rate 40%%, got %f", summary.SavingsRate)
	}
}

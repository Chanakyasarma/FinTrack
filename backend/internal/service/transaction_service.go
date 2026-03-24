package service

import (
	"context"
	"fmt"
	"time"

	"fintrack/internal/domain"
	"fintrack/internal/websocket"
)

type TransactionService struct {
	txRepo      domain.TransactionRepository
	accountRepo domain.AccountRepository
	cache       Cache
	hub         *websocket.Hub
}

func NewTransactionService(
	txRepo domain.TransactionRepository,
	accountRepo domain.AccountRepository,
	c Cache,
	hub *websocket.Hub,
) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		accountRepo: accountRepo,
		cache:       c,
		hub:         hub,
	}
}

// NewTransactionServiceForTest bypasses Redis for unit tests.
func NewTransactionServiceForTest(
	txRepo domain.TransactionRepository,
	accountRepo domain.AccountRepository,
	hub *websocket.Hub,
) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		accountRepo: accountRepo,
		cache:       &noopCacheImpl{},
		hub:         hub,
	}
}

type CreateTransactionInput struct {
	AccountID    string                    `json:"account_id"`
	Amount       float64                   `json:"amount"`
	Type         domain.TransactionType    `json:"type"`
	Category     domain.TransactionCategory `json:"category"`
	Description  string                    `json:"description"`
	MerchantName string                    `json:"merchant_name"`
}

func (s *TransactionService) Create(ctx context.Context, userID string, input CreateTransactionInput) (*domain.Transaction, error) {
	account, err := s.accountRepo.GetByID(ctx, input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	if account.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	tx := &domain.Transaction{
		AccountID:    input.AccountID,
		UserID:       userID,
		Amount:       input.Amount,
		Type:         input.Type,
		Category:     input.Category,
		Description:  input.Description,
		MerchantName: input.MerchantName,
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("creating transaction: %w", err)
	}

	// Update account balance
	newBalance := account.Balance
	if tx.Type == domain.TransactionTypeCredit {
		newBalance += tx.Amount
	} else {
		newBalance -= tx.Amount
	}

	if err := s.accountRepo.UpdateBalance(ctx, account.ID, newBalance); err != nil {
		return nil, fmt.Errorf("updating balance: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(ctx,
		txListKey(userID),
		summaryKey(userID),
		accountListKey(userID),
	)

	// Broadcast real-time event
	s.hub.SendToUser(userID, websocket.Event{
		Type: websocket.EventTransactionCreated,
		Payload: tx,
	})
	s.hub.SendToUser(userID, websocket.Event{
		Type: websocket.EventBalanceUpdated,
		Payload: map[string]any{
			"account_id": account.ID,
			"balance":    newBalance,
		},
	})

	return tx, nil
}

func (s *TransactionService) List(ctx context.Context, userID string, limit, offset int) ([]*domain.Transaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var txs []*domain.Transaction
	cacheKey := fmt.Sprintf("%s:%d:%d", txListKey(userID), limit, offset)

	if err := s.cache.Get(ctx, cacheKey, &txs); err == nil {
		return txs, nil
	}

	txs, err := s.txRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("listing transactions: %w", err)
	}

	s.cache.Set(ctx, cacheKey, txs, 2*time.Minute)
	return txs, nil
}

func (s *TransactionService) Summary(ctx context.Context, userID string) (*domain.Summary, error) {
	cacheKey := summaryKey(userID)

	var summary domain.Summary
	if err := s.cache.Get(ctx, cacheKey, &summary); err == nil {
		return &summary, nil
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	result, err := s.txRepo.GetSummary(ctx, userID, startOfMonth, now)
	if err != nil {
		return nil, fmt.Errorf("getting summary: %w", err)
	}

	trend, err := s.txRepo.GetMonthlyTrend(ctx, userID, 6)
	if err != nil {
		return nil, fmt.Errorf("getting trend: %w", err)
	}
	result.MonthlyTrend = trend

	if result.MonthlyIncome > 0 {
		result.SavingsRate = ((result.MonthlyIncome - result.MonthlyExpenses) / result.MonthlyIncome) * 100
	}

	s.cache.Set(ctx, cacheKey, result, 3*time.Minute)
	return result, nil
}

func txListKey(userID string) string {
	return fmt.Sprintf("transactions:list:%s", userID)
}

func summaryKey(userID string) string {
	return fmt.Sprintf("summary:%s", userID)
}

type UpdateTransactionInput struct {
	Amount       float64                    `json:"amount"`
	Type         domain.TransactionType     `json:"type"`
	Category     domain.TransactionCategory `json:"category"`
	Description  string                     `json:"description"`
	MerchantName string                     `json:"merchant_name"`
}

func (s *TransactionService) Update(ctx context.Context, userID, txID string, input UpdateTransactionInput) (*domain.Transaction, error) {
	// Verify ownership
	existing, err := s.txRepo.GetByID(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found")
	}
	if existing.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Reverse old balance effect, apply new one if amount/type changed
	if existing.Amount != input.Amount || existing.Type != input.Type {
		account, err := s.accountRepo.GetByID(ctx, existing.AccountID)
		if err != nil {
			return nil, fmt.Errorf("account not found: %w", err)
		}

		// Reverse old transaction
		newBalance := account.Balance
		if existing.Type == domain.TransactionTypeCredit {
			newBalance -= existing.Amount
		} else {
			newBalance += existing.Amount
		}

		// Apply new transaction
		if input.Type == domain.TransactionTypeCredit {
			newBalance += input.Amount
		} else {
			newBalance -= input.Amount
		}

		if err := s.accountRepo.UpdateBalance(ctx, account.ID, newBalance); err != nil {
			return nil, fmt.Errorf("updating balance: %w", err)
		}
	}

	updated := &domain.Transaction{
		ID:           txID,
		UserID:       userID,
		AccountID:    existing.AccountID,
		Amount:       input.Amount,
		Type:         input.Type,
		Category:     input.Category,
		Description:  input.Description,
		MerchantName: input.MerchantName,
		CreatedAt:    existing.CreatedAt,
	}

	if err := s.txRepo.Update(ctx, updated); err != nil {
		return nil, fmt.Errorf("updating transaction: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(ctx, txListKey(userID), summaryKey(userID), accountListKey(userID))

	return updated, nil
}

func (s *TransactionService) Delete(ctx context.Context, userID, txID string) error {
	existing, err := s.txRepo.GetByID(ctx, txID)
	if err != nil {
		return fmt.Errorf("transaction not found")
	}
	if existing.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Reverse the balance effect
	account, err := s.accountRepo.GetByID(ctx, existing.AccountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	newBalance := account.Balance
	if existing.Type == domain.TransactionTypeCredit {
		newBalance -= existing.Amount
	} else {
		newBalance += existing.Amount
	}

	if err := s.accountRepo.UpdateBalance(ctx, account.ID, newBalance); err != nil {
		return fmt.Errorf("updating balance: %w", err)
	}

	if err := s.txRepo.Delete(ctx, txID); err != nil {
		return fmt.Errorf("deleting transaction: %w", err)
	}

	s.cache.Delete(ctx, txListKey(userID), summaryKey(userID), accountListKey(userID))
	return nil
}

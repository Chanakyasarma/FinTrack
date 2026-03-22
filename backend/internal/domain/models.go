package domain

import (
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type AccountType string

const (
	AccountTypeChecking AccountType = "checking"
	AccountTypeSavings  AccountType = "savings"
	AccountTypeCredit   AccountType = "credit"
)

type Account struct {
	ID          string      `json:"id" db:"id"`
	UserID      string      `json:"user_id" db:"user_id"`
	Name        string      `json:"name" db:"name"`
	Type        AccountType `json:"type" db:"type"`
	Balance     float64     `json:"balance" db:"balance"`
	Currency    string      `json:"currency" db:"currency"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "credit"
	TransactionTypeDebit  TransactionType = "debit"
)

type TransactionCategory string

const (
	CategoryFood          TransactionCategory = "food"
	CategoryTransport     TransactionCategory = "transport"
	CategoryShopping      TransactionCategory = "shopping"
	CategoryEntertainment TransactionCategory = "entertainment"
	CategoryHealth        TransactionCategory = "health"
	CategorySalary        TransactionCategory = "salary"
	CategoryOther         TransactionCategory = "other"
)

type Transaction struct {
	ID          string              `json:"id" db:"id"`
	AccountID   string              `json:"account_id" db:"account_id"`
	UserID      string              `json:"user_id" db:"user_id"`
	Amount      float64             `json:"amount" db:"amount"`
	Type        TransactionType     `json:"type" db:"type"`
	Category    TransactionCategory `json:"category" db:"category"`
	Description string              `json:"description" db:"description"`
	MerchantName string             `json:"merchant_name" db:"merchant_name"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
}

type MonthlySummary struct {
	Month       string  `json:"month"`
	TotalCredit float64 `json:"total_credit"`
	TotalDebit  float64 `json:"total_debit"`
	NetFlow     float64 `json:"net_flow"`
}

type Summary struct {
	TotalBalance    float64          `json:"total_balance"`
	MonthlyIncome   float64          `json:"monthly_income"`
	MonthlyExpenses float64          `json:"monthly_expenses"`
	SavingsRate     float64          `json:"savings_rate"`
	CategoryBreakdown map[string]float64 `json:"category_breakdown"`
	MonthlyTrend    []MonthlySummary `json:"monthly_trend"`
}

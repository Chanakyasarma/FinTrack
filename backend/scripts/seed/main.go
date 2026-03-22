// scripts/seed/main.go
// Run with: go run scripts/seed/main.go
// Requires DATABASE_URL to be set (or a .env file).

package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"fintrack/internal/config"
	"fintrack/pkg/database"

	"golang.org/x/crypto/bcrypt"
)

type seedAccount struct {
	id       string
	name     string
	acctType string
	balance  float64
}

var merchants = map[string][]string{
	"food":          {"Starbucks", "Whole Foods", "Chipotle", "McDonald's", "Domino's", "Zomato"},
	"transport":     {"Uber", "Lyft", "Metro Card", "Shell Gas", "Ola"},
	"shopping":      {"Amazon", "Flipkart", "IKEA", "Nike", "H&M"},
	"entertainment": {"Netflix", "Spotify", "Steam", "BookMyShow", "Disney+"},
	"health":        {"CVS Pharmacy", "Apollo Clinic", "Gym Membership", "1mg"},
	"salary":        {"ACME Corp Payroll", "Freelance Transfer", "Consulting Fee"},
	"other":         {"ATM Withdrawal", "Bank Transfer", "Miscellaneous"},
}

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connecting to db: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("running migrations: %v", err)
	}

	ctx := context.Background()

	// ── Create demo user ─────────────────────────────────────────────────────
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userID string
	err = db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, full_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name
		RETURNING id
	`, "demo@fintrack.io", string(hash), "Demo User").Scan(&userID)
	if err != nil {
		log.Fatalf("creating user: %v", err)
	}
	log.Printf("✓ User: demo@fintrack.io / password123 (id=%s)", userID)

	// ── Create accounts ───────────────────────────────────────────────────────
	accounts := []seedAccount{
		{name: "Main Checking", acctType: "checking", balance: 4850.00},
		{name: "Savings Vault", acctType: "savings", balance: 12500.00},
		{name: "Credit Card", acctType: "credit", balance: -1200.00},
	}

	for i, a := range accounts {
		err = db.QueryRowContext(ctx, `
			INSERT INTO accounts (user_id, name, type, balance, currency)
			VALUES ($1, $2, $3, $4, 'USD')
			ON CONFLICT DO NOTHING
			RETURNING id
		`, userID, a.name, a.acctType, a.balance).Scan(&accounts[i].id)
		if err != nil {
			// Account may already exist from a previous seed run; query it
			db.QueryRowContext(ctx,
				`SELECT id FROM accounts WHERE user_id = $1 AND name = $2`,
				userID, a.name,
			).Scan(&accounts[i].id)
		}
		log.Printf("✓ Account: %s (id=%s)", a.name, accounts[i].id)
	}

	// ── Create transactions spanning 6 months ─────────────────────────────────
	rng := rand.New(rand.NewSource(42))
	now := time.Now()

	type txSeed struct {
		category string
		txType   string
		minAmt   float64
		maxAmt   float64
		freq     int // times per month
	}

	patterns := []txSeed{
		{"salary", "credit", 3000, 5000, 1},
		{"food", "debit", 5, 80, 15},
		{"transport", "debit", 10, 50, 10},
		{"shopping", "debit", 20, 300, 4},
		{"entertainment", "debit", 10, 50, 6},
		{"health", "debit", 15, 200, 2},
		{"other", "debit", 10, 100, 3},
	}

	txCount := 0
	for monthsBack := 5; monthsBack >= 0; monthsBack-- {
		monthStart := time.Date(now.Year(), now.Month()-time.Month(monthsBack), 1, 0, 0, 0, 0, now.Location())
		monthEnd := monthStart.AddDate(0, 1, 0)
		if monthEnd.After(now) {
			monthEnd = now
		}

		for _, p := range patterns {
			for i := 0; i < p.freq; i++ {
				// Random time within the month
				delta := monthEnd.Unix() - monthStart.Unix()
				txTime := time.Unix(monthStart.Unix()+rng.Int63n(delta), 0)

				amount := p.minAmt + rng.Float64()*(p.maxAmt-p.minAmt)
				amount = float64(int(amount*100)) / 100 // round to 2dp

				merchantList := merchants[p.category]
				merchant := merchantList[rng.Intn(len(merchantList))]

				acctIndex := 0 // checking for most things
				if p.category == "salary" {
					acctIndex = 0
				} else if p.txType == "credit" && rng.Float64() > 0.7 {
					acctIndex = 1 // sometimes into savings
				}
				if accounts[acctIndex].id == "" {
					continue
				}

				_, err := db.ExecContext(ctx, `
					INSERT INTO transactions
						(account_id, user_id, amount, type, category, description, merchant_name, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				`,
					accounts[acctIndex].id,
					userID,
					amount,
					p.txType,
					p.category,
					merchant+" purchase",
					merchant,
					txTime,
				)
				if err != nil {
					log.Printf("warn: inserting transaction: %v", err)
					continue
				}
				txCount++
			}
		}
	}

	log.Printf("✓ Seeded %d transactions across 6 months", txCount)
	log.Println("\n──────────────────────────────────────────")
	log.Println("  Demo credentials:")
	log.Println("  Email:    demo@fintrack.io")
	log.Println("  Password: password123")
	log.Println("──────────────────────────────────────────")
}

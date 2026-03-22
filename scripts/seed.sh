#!/usr/bin/env bash
# seed.sh — Populate FinTrack with realistic demo data
# Usage: ./scripts/seed.sh [API_BASE_URL]
# Default API_BASE_URL: http://localhost:8080

set -euo pipefail

API="${1:-http://localhost:8080}/api/v1"
EMAIL="demo@fintrack.dev"
PASSWORD="demo1234"
FULL_NAME="Alex Morgan"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}▶ $*${NC}"; }
warn() { echo -e "${YELLOW}⚠ $*${NC}"; }

# ── 1. Register user ──────────────────────────────────────────────────
log "Registering demo user..."
REGISTER=$(curl -sf -X POST "$API/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"full_name\":\"$FULL_NAME\"}" || true)

if echo "$REGISTER" | grep -q '"token"'; then
  TOKEN=$(echo "$REGISTER" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  log "Registered new user."
else
  warn "User may already exist — trying login..."
  LOGIN=$(curl -sf -X POST "$API/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
  TOKEN=$(echo "$LOGIN" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  log "Logged in as existing user."
fi

AUTH="Authorization: Bearer $TOKEN"

# ── 2. Create accounts ────────────────────────────────────────────────
log "Creating accounts..."

CHECKING=$(curl -sf -X POST "$API/accounts" \
  -H "Content-Type: application/json" -H "$AUTH" \
  -d '{"name":"Main Checking","type":"checking","currency":"USD"}')
CHECKING_ID=$(echo "$CHECKING" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
log "  ✓ Checking account: $CHECKING_ID"

SAVINGS=$(curl -sf -X POST "$API/accounts" \
  -H "Content-Type: application/json" -H "$AUTH" \
  -d '{"name":"High-Yield Savings","type":"savings","currency":"USD"}')
SAVINGS_ID=$(echo "$SAVINGS" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
log "  ✓ Savings account: $SAVINGS_ID"

CREDIT=$(curl -sf -X POST "$API/accounts" \
  -H "Content-Type: application/json" -H "$AUTH" \
  -d '{"name":"Visa Platinum","type":"credit","currency":"USD"}')
CREDIT_ID=$(echo "$CREDIT" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
log "  ✓ Credit account: $CREDIT_ID"

# ── 3. Seed transactions ──────────────────────────────────────────────
log "Seeding transactions..."

post_tx() {
  local account_id="$1" amount="$2" type="$3" category="$4" merchant="$5" desc="$6"
  curl -sf -X POST "$API/transactions" \
    -H "Content-Type: application/json" -H "$AUTH" \
    -d "{
      \"account_id\":\"$account_id\",
      \"amount\":$amount,
      \"type\":\"$type\",
      \"category\":\"$category\",
      \"merchant_name\":\"$merchant\",
      \"description\":\"$desc\"
    }" > /dev/null
}

# Income
post_tx "$CHECKING_ID"  4800   credit salary       "ACME Corp"        "Monthly salary"
post_tx "$SAVINGS_ID"   200    credit other        "Interest"         "Savings interest"

# Food & Dining
post_tx "$CREDIT_ID"    85.50  debit  food         "Whole Foods"      "Weekly groceries"
post_tx "$CREDIT_ID"    12.00  debit  food         "Starbucks"        "Morning coffee"
post_tx "$CREDIT_ID"    45.75  debit  food         "Chipotle"         "Team lunch"
post_tx "$CHECKING_ID"  22.00  debit  food         "Domino's"         "Pizza night"
post_tx "$CREDIT_ID"    67.30  debit  food         "Trader Joe's"     "Groceries"

# Transport
post_tx "$CHECKING_ID"  55.00  debit  transport    "Shell Gas"        "Fuel"
post_tx "$CREDIT_ID"    14.99  debit  transport    "Uber"             "Ride to airport"
post_tx "$CHECKING_ID"  120.00 debit  transport    "MTA Monthly"      "Subway pass"

# Shopping
post_tx "$CREDIT_ID"    249.99 debit  shopping     "Amazon"           "Electronics"
post_tx "$CREDIT_ID"    89.00  debit  shopping     "Nike"             "Running shoes"
post_tx "$CREDIT_ID"    35.00  debit  shopping     "H&M"              "Clothing"

# Entertainment
post_tx "$CREDIT_ID"    15.99  debit  entertainment "Netflix"         "Monthly subscription"
post_tx "$CREDIT_ID"    9.99   debit  entertainment "Spotify"         "Music subscription"
post_tx "$CHECKING_ID"  62.00  debit  entertainment "AMC Theaters"    "Movie + popcorn"

# Health
post_tx "$CHECKING_ID"  40.00  debit  health       "CVS Pharmacy"     "Prescriptions"
post_tx "$CREDIT_ID"    50.00  debit  health       "LA Fitness"       "Gym membership"
post_tx "$CHECKING_ID"  25.00  debit  health       "Co-pay"           "Doctor visit"

# Transfers to savings
post_tx "$SAVINGS_ID"   500.00 credit other        "Self Transfer"    "Monthly savings transfer"
post_tx "$CHECKING_ID"  500.00 debit  other        "Self Transfer"    "To savings"

log "  ✓ 20 transactions seeded"

# ── 4. Summary ────────────────────────────────────────────────────────
echo ""
log "Demo credentials:"
echo "  Email:    $EMAIL"
echo "  Password: $PASSWORD"
echo "  API URL:  $API"
echo ""
log "Open http://localhost:3000 and sign in!"

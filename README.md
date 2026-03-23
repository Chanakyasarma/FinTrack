# FinTrack — Real-Time Personal Finance Dashboard

A production-grade full-stack application built with **Go** and **React** demonstrating clean architecture, real-time WebSockets, JWT authentication, and Redis caching.

---

## Tech Stack

### Backend (Go)
| Layer | Technology |
|---|---|
| HTTP Router | `chi` — lightweight, idiomatic |
| Database | PostgreSQL 16 + `sqlx` |
| Cache | Redis 7 + `go-redis` |
| Auth | JWT (HS256) via `golang-jwt` |
| Real-time | WebSockets via `gorilla/websocket` |
| Password | `bcrypt` |

### Frontend (React + TypeScript)
| Concern | Technology |
|---|---|
| Framework | React 18 + TypeScript |
| Routing | React Router v6 |
| Styling | Tailwind CSS |
| Charts | Recharts |
| HTTP | Axios with interceptors |
| Notifications | react-hot-toast |

### Infrastructure
- **Docker Compose** — one-command local setup
- **Nginx** — SPA serving + reverse proxy
- **Multi-stage Docker builds** — minimal production images

---

## Architecture

```
fintrack/
├── backend/
│   ├── main.go                         # Entry point, DI wiring, graceful shutdown
│   ├── internal/
│   │   ├── config/        config.go    # Environment config
│   │   ├── domain/        models.go    # Domain models (User, Account, Transaction)
│   │   │                  repository.go # Repository interfaces (ports)
│   │   ├── handler/       auth.go      # HTTP handlers (thin layer)
│   │   │                  account.go
│   │   │                  transaction.go
│   │   │                  websocket.go
│   │   ├── service/       auth.go      # Business logic
│   │   │                  account.go
│   │   │                  transaction.go
│   │   ├── repository/    user.go      # DB implementations
│   │   │                  account.go
│   │   │                  transaction.go
│   │   └── websocket/     hub.go       # WS hub, client pump goroutines
│   └── pkg/
│       ├── cache/         redis.go     # Redis client wrapper
│       ├── database/      postgres.go  # DB connection + migrations
│       ├── middleware/    jwt.go       # JWT auth middleware
│       └── respond/       respond.go   # HTTP response helpers
│
├── frontend/
│   └── src/
│       ├── context/       AuthContext.tsx   # JWT auth state
│       ├── hooks/         useFinance.ts     # Data fetching hooks
│       │                  useWebSocket.ts   # WS connection + auto-reconnect
│       ├── components/    StatCard.tsx
│       │                  TransactionItem.tsx
│       │                  AccountCard.tsx
│       │                  AddTransactionModal.tsx
│       │                  AddAccountModal.tsx
│       │                  Sidebar.tsx
│       ├── pages/         Dashboard.tsx
│       │                  TransactionsPage.tsx
│       │                  AccountsPage.tsx
│       │                  AnalyticsPage.tsx
│       │                  AuthPage.tsx
│       ├── lib/           api.ts       # Axios instance + interceptors
│       │                  format.ts    # Currency, date helpers
│       └── types/         index.ts     # Shared TypeScript types
│
└── docker-compose.yml
```

### Clean Architecture (Backend)

```
Request → Handler → Service → Repository → Database
                ↑         ↑
           (validates) (business logic,
                        cache, WS events)
```

- **Handlers** — decode request, validate input, call service, write response
- **Services** — own business logic, call repository, invalidate cache, emit WS events
- **Repositories** — implement domain interfaces, execute SQL
- **Domain** — pure models + repository interfaces; no framework dependencies

---

## API Endpoints

### Auth
```
POST /api/v1/auth/register   { email, password, full_name }
POST /api/v1/auth/login      { email, password }
```

### Accounts (JWT required)
```
GET  /api/v1/accounts
POST /api/v1/accounts        { name, type, currency }
GET  /api/v1/accounts/:id
```

### Transactions (JWT required)
```
GET  /api/v1/transactions?limit=20&offset=0
POST /api/v1/transactions    { account_id, amount, type, category, description, merchant_name }
GET  /api/v1/transactions/summary
```

### WebSocket (JWT required)
```
WS   /api/v1/ws              Real-time event stream
```

**WS Event types:**
```json
{ "type": "transaction.created", "payload": { ...transaction } }
{ "type": "balance.updated",     "payload": { "account_id": "...", "balance": 1234.56 } }
```

---

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for local dev)
- Node.js 20+ (for local dev)

### With Docker (recommended)
```bash
git clone <repo>
cd fintrack

cp backend/.env.example backend/.env
# Edit JWT_SECRET in backend/.env

docker-compose up --build
```

App: http://localhost:3000  
API: http://localhost:8080

### Local Development

**Backend:**
```bash
cd backend
cp .env.example .env
go mod download
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

---

## Key Design Decisions

### Repository Pattern with Interfaces
All repositories are defined as interfaces in `domain/repository.go`. This makes them trivially mockable for unit tests and decouples business logic from the database driver.

```go
// Service depends on the interface, not the concrete type
type TransactionService struct {
    txRepo domain.TransactionRepository  // ← interface
    ...
}
```

### Cache Invalidation Strategy
Uses a **write-through invalidation** pattern: on any write (new transaction, balance update), affected cache keys are deleted. On next read, the DB is queried and the result is re-cached with a short TTL.

### WebSocket Hub
A single `Hub` goroutine serialises all register/unregister/broadcast operations, avoiding lock contention. Each user's clients are tracked in a `map[string]map[*Client]struct{}` keyed by `userID`, enabling targeted per-user broadcasting.

### Graceful Shutdown
Backend listens for `SIGINT`/`SIGTERM` and calls `http.Server.Shutdown()` with a 30-second context, allowing in-flight requests to complete.

---

## Interview Talking Points

| Topic | What to say |
|---|---|
| **Architecture** | "I used Clean Architecture — handler/service/repository layers. Domain defines interfaces; concrete implementations live in separate packages. This makes the codebase testable and easy to extend." |
| **WebSockets** | "The Hub pattern uses a single goroutine for channel operations to avoid race conditions. Each client has send/receive pump goroutines. I support per-user broadcasting keyed by userID." |
| **Caching** | "Redis caches account lists (5min TTL) and summary data (3min TTL). Any write operation invalidates affected keys. I chose short TTLs over event-driven cache updates to keep it simple without a message queue." |
| **Auth** | "JWT HS256 with 72-hour expiry. The middleware extracts the userID into the request context. Passwords are bcrypt-hashed. In production I'd add refresh tokens and token revocation via Redis." |
| **SQL** | "I wrote raw SQL with sqlx for full control. Indexes on user_id, account_id, and created_at DESC cover the most common query patterns." |
| **Error handling** | "Errors are wrapped with context using fmt.Errorf. Domain errors (ErrInvalidCredentials) are defined as sentinel errors and checked with errors.Is for clean HTTP status mapping." |

---

## Potential Extensions (discuss in interviews)

- **Pagination cursor** — replace limit/offset with cursor-based pagination for large datasets
- **Refresh tokens** — short-lived access tokens + long-lived refresh tokens stored in Redis
- **Rate limiting** — middleware using Redis sliding window counters
- **Budget goals** — alert via WebSocket when a spending category exceeds a threshold
- **CSV export** — streaming export of transactions using `encoding/csv`
- **Unit tests** — mock repositories with `testify/mock`, table-driven tests for services
- **OpenTelemetry** — distributed tracing across handler/service/repo layers
<img width="1918" height="992" alt="Screenshot 1" src="https://github.com/user-attachments/assets/beb07339-e2bc-492b-9570-291608bee8de" />
<img width="1917" height="988" alt="Screenshot 2026-03-23 135019" src="https://github.com/user-attachments/assets/9a8dfc73-89f1-48fa-b195-d2e3f2748729" />



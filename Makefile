.PHONY: help dev up down build test lint seed clean tidy

# ── Colors ────────────────────────────────────────────────────────────────────
CYAN  := \033[0;36m
RESET := \033[0m

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-18s$(RESET) %s\n", $$1, $$2}'

# ── Docker ────────────────────────────────────────────────────────────────────
up: ## Start all services (Postgres, Redis, Backend, Frontend)
	docker-compose up --build -d
	@echo "$(CYAN)→ App running at http://localhost:3000$(RESET)"
	@echo "$(CYAN)→ API running at http://localhost:8080$(RESET)"

down: ## Stop all services
	docker-compose down

logs: ## Tail logs for all services
	docker-compose logs -f

# ── Local Dev ─────────────────────────────────────────────────────────────────
infra: ## Start only Postgres + Redis (for local backend dev)
	docker-compose up -d postgres redis
	@echo "$(CYAN)→ Postgres: localhost:5432$(RESET)"
	@echo "$(CYAN)→ Redis:    localhost:6379$(RESET)"

backend: ## Run backend locally (requires infra running)
	cd backend && go run main.go

frontend: ## Run frontend dev server
	cd frontend && npm run dev

dev: ## Start infra + backend + frontend (3 terminals needed, or use tmux)
	$(MAKE) infra
	@echo "Run in separate terminals:"
	@echo "  make backend"
	@echo "  make frontend"

# ── Go ────────────────────────────────────────────────────────────────────────
tidy: ## go mod tidy
	cd backend && go mod tidy

build: ## Build backend binary
	cd backend && go build -ldflags="-s -w" -o bin/fintrack ./main.go
	@echo "$(CYAN)→ Binary: backend/bin/fintrack$(RESET)"

test: ## Run all backend tests
	cd backend && go test ./... -v -count=1

test-cover: ## Run tests with coverage report
	cd backend && go test ./... -coverprofile=coverage.out
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "$(CYAN)→ Coverage report: backend/coverage.html$(RESET)"

lint: ## Run golangci-lint (install: https://golangci-lint.run)
	cd backend && golangci-lint run ./...

vet: ## Run go vet
	cd backend && go vet ./...

# ── Frontend ──────────────────────────────────────────────────────────────────
fe-install: ## Install frontend dependencies
	cd frontend && npm install

fe-build: ## Build frontend for production
	cd frontend && npm run build

fe-check: ## TypeScript type check
	cd frontend && npx tsc --noEmit

# ── Database ──────────────────────────────────────────────────────────────────
seed: ## Seed the database with sample data
	cd backend && go run scripts/seed/main.go

psql: ## Open a psql shell to the running Postgres container
	docker-compose exec postgres psql -U postgres -d fintrack

redis-cli: ## Open redis-cli to the running Redis container
	docker-compose exec redis redis-cli

# ── Misc ──────────────────────────────────────────────────────────────────────
clean: ## Remove build artifacts
	rm -rf backend/bin backend/coverage.out backend/coverage.html frontend/dist

env: ## Copy .env.example to .env (if not already present)
	@[ -f backend/.env ] || cp backend/.env.example backend/.env && echo "$(CYAN)→ Created backend/.env$(RESET)"

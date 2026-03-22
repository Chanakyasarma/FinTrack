# Contributing to FinTrack

## Project Philosophy

This codebase prioritises **clarity over cleverness**. Every pattern used should
be explainable in a 5-minute interview conversation.

---

## Code Conventions

### Go — Backend

#### Error Handling
Always wrap errors with context using `fmt.Errorf`:

```go
// ✅ Good — the caller knows where the error came from
if err := s.repo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("creating user: %w", err)
}

// ❌ Bad — loses callsite context
return nil, err
```

Sentinel errors live in the service layer and are checked with `errors.Is`:

```go
// Define in service
var ErrNotFound = errors.New("not found")

// Check in handler
if errors.Is(err, service.ErrNotFound) {
    respond.Error(w, http.StatusNotFound, "resource not found")
    return
}
```

#### Context propagation
Every function that touches I/O (DB, Redis, HTTP) must accept a `context.Context`
as its first parameter. Never use `context.Background()` inside a handler — pass
the request context from `r.Context()`.

#### Naming
- Interfaces: describe behaviour, not implementation (`UserRepository` not `UserStore`)
- Constructors: `New<TypeName>(deps...) *TypeName`
- Test constructors: `New<TypeName>ForTest(...)` or `New<TypeName>WithNilCache(...)`

#### Package structure
```
internal/  — application code; not importable externally
pkg/       — general utilities (cache, db, middleware); could be extracted to a library
```

#### Table-driven tests
Prefer table-driven tests for input/output variations:

```go
tests := []struct {
    name     string
    input    service.RegisterInput
    wantErr  error
}{
    {"valid", service.RegisterInput{...}, nil},
    {"duplicate", service.RegisterInput{...}, service.ErrUserAlreadyExists},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        _, err := svc.Register(ctx, tt.input)
        if !errors.Is(err, tt.wantErr) {
            t.Errorf("got %v, want %v", err, tt.wantErr)
        }
    })
}
```

---

### React — Frontend

#### Custom hooks for all data fetching
Never call `api.get(...)` inside a component directly. Use a custom hook:

```ts
// ✅ Good
const { accounts, loading } = useAccounts()

// ❌ Bad — couples component to fetch logic
const [accounts, setAccounts] = useState([])
useEffect(() => { api.get('/accounts').then(...) }, [])
```

#### Component sizing
Components that exceed ~150 lines should be split. Signs a component is too large:
- More than one modal
- More than one `useEffect`
- More than two pieces of local state

#### Type discipline
Every API response shape must have a corresponding type in `src/types/index.ts`.
Never use `any` — if the shape is unknown, use `unknown` and narrow it explicitly.

---

## Running Tests

```bash
# All backend tests
make test

# With coverage
make test-cover

# Single package
cd backend && go test ./internal/service/... -v

# TypeScript check
make fe-check
```

## Adding a New Feature (checklist)

- [ ] Add domain model to `internal/domain/models.go`
- [ ] Add repository interface method to `internal/domain/repository.go`
- [ ] Implement in `internal/repository/`
- [ ] Write service logic in `internal/service/` (cache + WS where appropriate)
- [ ] Add handler in `internal/handler/`
- [ ] Register route in `main.go`
- [ ] Add TypeScript type in `frontend/src/types/index.ts`
- [ ] Add API call in `frontend/src/hooks/`
- [ ] Build UI component
- [ ] Write unit tests for service + handler

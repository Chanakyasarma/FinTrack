package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fintrack/internal/domain"
	"fintrack/internal/handler"
	"fintrack/internal/service"
)

// --- Minimal in-memory user repo ---

type inMemoryUserRepo struct {
	users map[string]*domain.User
	byEmail map[string]*domain.User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{
		users:   make(map[string]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (r *inMemoryUserRepo) Create(ctx context.Context, u *domain.User) error {
	if _, exists := r.byEmail[u.Email]; exists {
		return errors.New("duplicate")
	}
	u.ID = "user-id-1"
	r.users[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *inMemoryUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := r.byEmail[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (r *inMemoryUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

// --- Helpers ---

func newAuthHandler() *handler.AuthHandler {
	repo := newInMemoryUserRepo()
	svc := service.NewAuthService(repo, "test-secret-key")
	return handler.NewAuthHandler(svc)
}

func postJSON(handler http.HandlerFunc, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

// --- Tests ---

func TestAuthHandler_Register_Success(t *testing.T) {
	h := newAuthHandler()

	rr := postJSON(h.Register, "/api/v1/auth/register", map[string]string{
		"email":     "test@example.com",
		"password":  "password123",
		"full_name": "Test User",
	})

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data, _ := resp["data"].(map[string]any)
	if data["token"] == nil || data["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	h := newAuthHandler()

	rr := postJSON(h.Register, "/api/v1/auth/register", map[string]string{
		"email": "test@example.com",
		// missing password and full_name
	})

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
}

func TestAuthHandler_Register_Duplicate(t *testing.T) {
	h := newAuthHandler()

	body := map[string]string{
		"email":     "dup@example.com",
		"password":  "password123",
		"full_name": "Dup User",
	}
	postJSON(h.Register, "/api/v1/auth/register", body)
	rr := postJSON(h.Register, "/api/v1/auth/register", body)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d", rr.Code)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	h := newAuthHandler()

	postJSON(h.Register, "/api/v1/auth/register", map[string]string{
		"email":     "login@example.com",
		"password":  "mypassword",
		"full_name": "Login User",
	})

	rr := postJSON(h.Login, "/api/v1/auth/login", map[string]string{
		"email":    "login@example.com",
		"password": "mypassword",
	})

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	h := newAuthHandler()

	postJSON(h.Register, "/api/v1/auth/register", map[string]string{
		"email":     "wrong@example.com",
		"password":  "correctpass",
		"full_name": "Wrong Pass",
	})

	rr := postJSON(h.Login, "/api/v1/auth/login", map[string]string{
		"email":    "wrong@example.com",
		"password": "wrongpass",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthHandler_Login_BadBody(t *testing.T) {
	h := newAuthHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

package service_test

import (
	"context"
	"errors"
	"testing"

	"fintrack/internal/domain"
	"fintrack/internal/service"
)

// --- Mock UserRepository ---

type mockUserRepo struct {
	users  map[string]*domain.User
	byEmail map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[string]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	if _, exists := m.byEmail[user.Email]; exists {
		return errors.New("duplicate email")
	}
	user.ID = "test-user-id"
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

// --- Tests ---

func TestRegister_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewAuthService(repo, "test-secret")

	result, err := svc.Register(context.Background(), service.RegisterInput{
		Email:    "alice@example.com",
		Password: "password123",
		FullName: "Alice Smith",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected non-empty token")
	}
	if result.User.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", result.User.Email)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewAuthService(repo, "test-secret")

	input := service.RegisterInput{
		Email:    "alice@example.com",
		Password: "password123",
		FullName: "Alice Smith",
	}

	if _, err := svc.Register(context.Background(), input); err != nil {
		t.Fatalf("first registration failed unexpectedly: %v", err)
	}

	_, err := svc.Register(context.Background(), input)
	if !errors.Is(err, service.ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got: %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewAuthService(repo, "test-secret")

	_, err := svc.Register(context.Background(), service.RegisterInput{
		Email:    "bob@example.com",
		Password: "securepass",
		FullName: "Bob Jones",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	result, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "bob@example.com",
		Password: "securepass",
	})
	if err != nil {
		t.Fatalf("expected login to succeed, got: %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected non-empty JWT token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewAuthService(repo, "test-secret")

	_, _ = svc.Register(context.Background(), service.RegisterInput{
		Email:    "carol@example.com",
		Password: "correctpassword",
		FullName: "Carol White",
	})

	_, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "carol@example.com",
		Password: "wrongpassword",
	})
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := service.NewAuthService(repo, "test-secret")

	_, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "nobody@example.com",
		Password: "anypassword",
	})
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cinema/auth-service/internal/domain"
	"cinema/auth-service/internal/usecase"
	"github.com/google/uuid"
)


type mockUserRepo struct {
	users map[uuid.UUID]*domain.User
}
func newMockUserRepo() *mockUserRepo { return &mockUserRepo{users: make(map[uuid.UUID]*domain.User)} }
func (m *mockUserRepo) Create(_ context.Context, u *domain.User) error { m.users[u.ID] = u; return nil }
func (m *mockUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok { return nil, domain.ErrNotFound }
	return u, nil
}
func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email { return u, nil }
	}
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Update(_ context.Context, u *domain.User) error { m.users[u.ID] = u; return nil }
func (m *mockUserRepo) Delete(_ context.Context, id uuid.UUID) error { delete(m.users, id); return nil }

type mockTokenCache struct {
	blacklist     map[string]bool
	refreshTokens map[string]string
}
func newMockTokenCache() *mockTokenCache {
	return &mockTokenCache{
		blacklist:     make(map[string]bool),
		refreshTokens: make(map[string]string),
	}
}
func (c *mockTokenCache) BlacklistToken(_ context.Context, token string, _ time.Duration) error {
	c.blacklist[token] = true; return nil
}
func (c *mockTokenCache) IsBlacklisted(_ context.Context, token string) (bool, error) {
	return c.blacklist[token], nil
}
func (c *mockTokenCache) StoreRefreshToken(_ context.Context, userID, token string, _ time.Duration) error {
	c.refreshTokens[token] = userID; return nil
}
func (c *mockTokenCache) GetRefreshToken(_ context.Context, token string) (string, error) {
	v, ok := c.refreshTokens[token]
	if !ok { return "", errors.New("not found") }
	return v, nil
}
func (c *mockTokenCache) DeleteRefreshToken(_ context.Context, token string) error {
	delete(c.refreshTokens, token); return nil
}

type mockEmail struct{ sent []string }
func (e *mockEmail) SendWelcome(email, name string) error { e.sent = append(e.sent, email); return nil }
func (e *mockEmail) SendBookingConfirmed(_, _, _, _ string, _ []string, _ float64) error { return nil }

func newTestUsecase() (*usecase.AuthUsecase, *mockUserRepo, *mockTokenCache, *mockEmail) {
	repo  := newMockUserRepo()
	cache := newMockTokenCache()
	email := &mockEmail{}
	uc    := usecase.NewAuthUsecase(repo, cache, email)
	return uc, repo, cache, email
}


func TestRegister_Success(t *testing.T) {
	uc, repo, _, _ := newTestUsecase()
	user, err := uc.Register(context.Background(), "test@example.com", "password123", "Test User", "123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("wrong email: %s", user.Email)
	}
	if _, ok := repo.users[user.ID]; !ok {
		t.Error("user not stored in repo")
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	uc, _, _, _ := newTestUsecase()
	_, err := uc.Register(context.Background(), "test@example.com", "short", "Test", "")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	uc, _, _, _ := newTestUsecase()
	uc.Register(context.Background(), "dupe@example.com", "password123", "User One", "")
	_, err := uc.Register(context.Background(), "dupe@example.com", "password456", "User Two", "")
	if !errors.Is(err, domain.ErrEmailExists) {
		t.Errorf("expected ErrEmailExists, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	uc, _, cache, _ := newTestUsecase()
	uc.Register(context.Background(), "login@example.com", "mypassword", "Login User", "")
	pair, user, err := uc.Login(context.Background(), "login@example.com", "mypassword")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if pair.AccessToken == "" {
		t.Error("empty access token")
	}
	if user.Email != "login@example.com" {
		t.Error("wrong user returned")
	}
	if len(cache.refreshTokens) != 1 {
		t.Error("refresh token not stored")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	uc, _, _, _ := newTestUsecase()
	uc.Register(context.Background(), "a@b.com", "correctpass", "User", "")
	_, _, err := uc.Login(context.Background(), "a@b.com", "wrongpass")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogout_BlacklistsToken(t *testing.T) {
	uc, _, cache, _ := newTestUsecase()
	uc.Register(context.Background(), "out@example.com", "password123", "User", "")
	pair, _, _ := uc.Login(context.Background(), "out@example.com", "password123")

	err := uc.Logout(context.Background(), pair.AccessToken)
	if err != nil {
		t.Fatalf("logout failed: %v", err)
	}
	if !cache.blacklist[pair.AccessToken] {
		t.Error("token not blacklisted")
	}
}

func TestValidateToken_Blacklisted(t *testing.T) {
	uc, _, _, _ := newTestUsecase()
	uc.Register(context.Background(), "val@example.com", "password123", "User", "")
	pair, _, _ := uc.Login(context.Background(), "val@example.com", "password123")
	uc.Logout(context.Background(), pair.AccessToken)

	_, err := uc.ValidateToken(context.Background(), pair.AccessToken)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

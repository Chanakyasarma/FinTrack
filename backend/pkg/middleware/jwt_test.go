package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fintrack/pkg/middleware"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key"

func makeToken(userID string, secret string, exp time.Time) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func nextHandler(t *testing.T, expectUserID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gotID := middleware.GetUserID(r)
		if gotID != expectUserID {
			t.Errorf("expected userID %q in context, got %q", expectUserID, gotID)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	token := makeToken("user-abc", testSecret, time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	mw := middleware.JWTAuth(testSecret)(nextHandler(t, "user-abc"))
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()

	mw := middleware.JWTAuth(testSecret)(nextHandler(t, ""))
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	token := makeToken("user-expired", testSecret, time.Now().Add(-time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	mw := middleware.JWTAuth(testSecret)(nextHandler(t, ""))
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for expired token, got %d", rr.Code)
	}
}

func TestJWTAuth_WrongSecret(t *testing.T) {
	token := makeToken("user-xyz", "wrong-secret", time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	mw := middleware.JWTAuth(testSecret)(nextHandler(t, ""))
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong secret, got %d", rr.Code)
	}
}

func TestJWTAuth_MalformedHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token not-a-bearer-format")
	rr := httptest.NewRecorder()

	mw := middleware.JWTAuth(testSecret)(nextHandler(t, ""))
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for malformed header, got %d", rr.Code)
	}
}

package authhandler

import (
	"net/http"
	"net/http/httptest"
	"picup/core/infrastructure/config"
	"strings"
	"testing"
)

func TestRequireAuth_NoHeader(t *testing.T) {
	config.Server.Auth.JWTSecret = "secret"
	config.Server.Auth.JWTExpirationSeconds = 3600

	mw := NewAuthMiddleware()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/images", nil)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	config.Server.Auth.JWTSecret = "secret"
	config.Server.Auth.JWTExpirationSeconds = 3600

	mw := NewAuthMiddleware()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/images", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	config.Server.Auth.AdminUsername = "admin"
	config.Server.Auth.AdminPassword = "pass"
	config.Server.Auth.JWTSecret = "secret"
	config.Server.Auth.JWTExpirationSeconds = 3600

	h := NewAuthHandler()

	body := `{"username":"admin","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestMe_Unauthorized_NoCookie(t *testing.T) {
	config.Server.Auth.JWTSecret = "secret"
	config.Server.Auth.JWTExpirationSeconds = 3600

	h := NewAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMe_Unauthorized_InvalidCookie(t *testing.T) {
	config.Server.Auth.JWTSecret = "secret"
	config.Server.Auth.JWTExpirationSeconds = 3600

	h := NewAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "invalid"})
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMe_Success_ValidCookie(t *testing.T) {
	config.Server.Auth.JWTSecret = "secret"
	config.Server.Auth.JWTExpirationSeconds = 3600
	config.Server.Auth.AdminUsername = "admin"
	config.Server.Auth.AdminPassword = "pass"

	h := NewAuthHandler()

	// Generate a valid token by calling Login.
	body := `{"username":"admin","password":"pass"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	h.Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", loginW.Code)
	}

	// Extract cookie from Set-Cookie header.
	setCookie := loginW.Header().Get("Set-Cookie")
	if setCookie == "" {
		t.Fatalf("expected Set-Cookie header")
	}

	// Very small parser: cookie value is before first ';' and after name '='.
	parts := strings.SplitN(setCookie, ";", 2)
	nameValue := parts[0]
	nameParts := strings.SplitN(nameValue, "=", 2)
	if len(nameParts) != 2 {
		t.Fatalf("unexpected Set-Cookie format: %s", setCookie)
	}
	cookieValue := nameParts[1]

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: cookieValue})
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

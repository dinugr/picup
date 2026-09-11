package authhandler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticImagesRouteNotProtected(t *testing.T) {
	// This test is a routing-level regression: auth middleware should only be applied
	// to /api/images* endpoints, not to /images/* static file routes.
	mux := http.NewServeMux()

	// Simulate static route registration (no auth middleware).
	mux.Handle("/images/master/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Simulate protected route registration.
	mw := NewAuthMiddleware()
	mux.Handle("/api/images", mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	// Static should be reachable without Authorization.
	staticReq := httptest.NewRequest(http.MethodGet, "/images/master/hello", nil)
	staticRes := httptest.NewRecorder()
	mux.ServeHTTP(staticRes, staticReq)
	if staticRes.Code != http.StatusOK {
		t.Fatalf("expected static route 200, got %d", staticRes.Code)
	}

	// Protected should be rejected without Authorization.
	protectedReq := httptest.NewRequest(http.MethodGet, "/api/images", nil)
	protectedRes := httptest.NewRecorder()
	mux.ServeHTTP(protectedRes, protectedReq)
	if protectedRes.Code != http.StatusUnauthorized {
		t.Fatalf("expected protected route 401, got %d", protectedRes.Code)
	}
}

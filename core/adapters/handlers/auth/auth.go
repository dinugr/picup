package authhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"picup/core/domain/models"
	"picup/core/domain/repositories"
	"picup/core/infrastructure/config"
	"picup/core/infrastructure/security"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	jwtCfg security.JWTConfig
	repo   repositories.Repository
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken      string `json:"accessToken"`
	TokenType        string `json:"tokenType"`
	ExpiresInSeconds int64  `json:"expiresInSeconds"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewAuthHandler(repo ...repositories.Repository) *AuthHandler {
	jwtCfg := security.NewJWTConfig(
		config.Server.Auth.JWTSecret,
		config.Server.Auth.JWTIssuer,
		config.Server.Auth.JWTAudience,
		time.Duration(config.Server.Auth.JWTExpirationSeconds)*time.Second,
	)
	var sessionRepo repositories.Repository
	if len(repo) > 0 {
		sessionRepo = repo[0]
	}
	return &AuthHandler{jwtCfg: jwtCfg, repo: sessionRepo}
}

func (h *AuthHandler) RouteInit(mux *http.ServeMux, chain func(http.Handler) http.Handler) {
	mux.Handle("POST /api/auth/login", chain(http.HandlerFunc(h.Login)))
	mux.Handle("GET /api/auth/session", chain(http.HandlerFunc(h.Me)))
	mux.Handle("POST /api/auth/logout", chain(http.HandlerFunc(h.Logout)))
	mux.Handle("GET /api/auth/sessions", chain(http.HandlerFunc(h.ListSessions)))
	mux.Handle("DELETE /api/auth/sessions/{sessionID}", chain(http.HandlerFunc(h.RevokeSession)))
	mux.Handle("POST /api/auth/logout-all", chain(http.HandlerFunc(h.LogoutAll)))
}

// Me validates the HttpOnly session cookie and returns 200 when authenticated.
// Frontend uses this to rehydrate auth state on refresh.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
		return
	}

	if strings.TrimSpace(cookie.Value) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
		return
	}

	claims, err := security.ValidateAccessToken(h.jwtCfg, cookie.Value)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
		return
	}
	if h.repo != nil {
		if session, lookupErr := h.repo.GetSessionByJTI(r.Context(), claims.ID); lookupErr != nil || !session.IsActive(time.Now()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"authenticated":true}`))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_request"})
		return
	}

	username := strings.TrimSpace(req.Username)
	password := req.Password

	if username == "" || password == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid_credentials"})
		return
	}

	if username != config.Server.Auth.AdminUsername || password != config.Server.Auth.AdminPassword {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid_credentials"})
		return
	}

	token, err := security.GenerateAccessToken(h.jwtCfg, username)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "token_generation_failed"})
		log.Printf("[WARNING] %s", err.Error())
		return
	}
	claims, err := security.ValidateAccessToken(h.jwtCfg, token)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "token_generation_failed"})
		return
	}
	if h.repo != nil {
		now := time.Now().UTC()
		if err := h.repo.CreateSession(r.Context(), &models.Session{
			ID: uuid.NewString(), UserID: username, TokenJTI: claims.ID,
			CreatedAt: now, ExpiresAt: now.Add(h.jwtCfg.Expiration),
			UserAgent: r.UserAgent(), IPAddress: clientIP(r),
		}); err != nil {
			h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "session_persistence_failed"})
			log.Printf("[WARNING] session persistence failed: %v", err)
			return
		}
	}

	expiresIn := int64(h.jwtCfg.Expiration.Seconds())
	if expiresIn <= 0 {
		expiresIn = 3600
	}

	// Set JWT as HttpOnly cookie so frontend cannot read it (mitigates XSS token theft).
	// Important: `Secure` cookies are only sent over HTTPS.
	// Decide based on the actual request (or proxy forwarded proto), not config host.
	cookieSecure := r.TLS != nil
	if !cookieSecure {
		if strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
			cookieSecure = true
		}
	}

	cookieParts := []string{
		sessionCookieName + "=" + token,
		"Path=/",
		"HttpOnly",
		"SameSite=Lax",
		"Max-Age=" + fmt.Sprintf("%d", expiresIn),
	}
	if cookieSecure {
		cookieParts = append(cookieParts, "Secure")
	}

	w.Header().Set("Set-Cookie", strings.Join(cookieParts, "; "))

	// Keep response body for compatibility with existing frontend.
	h.writeJSON(w, http.StatusOK, loginResponse{
		AccessToken:      token,
		TokenType:        "Bearer",
		ExpiresInSeconds: expiresIn,
	})
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		// best-effort; nothing else we can do
		_ = fmt.Errorf("failed to encode json: %v", err)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, userID, ok := h.authenticatedClaims(r)
	if !ok {
		writeUnauthorized(w)
		return
	}
	if err := h.repo.RevokeSessionByJTI(r.Context(), claims.ID, userID, "logout", time.Now()); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "logout_failed"})
		return
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	_, userID, ok := h.authenticatedClaims(r)
	if !ok {
		writeUnauthorized(w)
		return
	}
	sessions, err := h.repo.ListActiveSessions(r.Context(), userID, time.Now())
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "session_list_failed"})
		return
	}
	currentClaims, _, _ := h.authenticatedClaims(r)
	type sessionResponse struct {
		ID         string     `json:"id"`
		UserAgent  string     `json:"userAgent"`
		CreatedAt  time.Time  `json:"createdAt"`
		LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
		Current    bool       `json:"current"`
	}
	result := make([]sessionResponse, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, sessionResponse{ID: session.ID, UserAgent: session.UserAgent, CreatedAt: session.CreatedAt, LastSeenAt: session.LastSeenAt, Current: currentClaims != nil && session.TokenJTI == currentClaims.ID})
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *AuthHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	_, userID, ok := h.authenticatedClaims(r)
	if !ok {
		writeUnauthorized(w)
		return
	}
	if err := h.repo.RevokeSession(r.Context(), r.PathValue("sessionID"), userID, "logout", time.Now()); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "session_revoke_failed"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	_, userID, ok := h.authenticatedClaims(r)
	if !ok {
		writeUnauthorized(w)
		return
	}
	if err := h.repo.RevokeAllSessions(r.Context(), userID, "logout_all", time.Now()); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "logout_all_failed"})
		return
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) authenticatedClaims(r *http.Request) (*security.Claims, string, bool) {
	if h.repo == nil {
		return nil, "", false
	}
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return nil, "", false
	}
	claims, err := security.ValidateAccessToken(h.jwtCfg, cookie.Value)
	if err != nil {
		return nil, "", false
	}
	session, err := h.repo.GetSessionByJTI(r.Context(), claims.ID)
	if err != nil || !session.IsActive(time.Now()) {
		return nil, "", false
	}
	return claims, session.UserID, true
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: sessionCookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); forwarded != "" {
		return forwarded
	}
	if host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

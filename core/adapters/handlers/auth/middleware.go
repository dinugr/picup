package authhandler

import (
	"net/http"
	"picup/core/domain/repositories"
	"picup/core/infrastructure/config"
	"picup/core/infrastructure/security"
	"strings"
	"time"
)

const sessionCookieName = "picup-session"

type AuthMiddleware struct {
	jwtCfg security.JWTConfig
	repo   repositories.Repository
}

func NewAuthMiddleware(repo ...repositories.Repository) *AuthMiddleware {
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
	return &AuthMiddleware{jwtCfg: jwtCfg, repo: sessionRepo}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			writeUnauthorized(w)
			return
		}
		if strings.TrimSpace(cookie.Value) == "" {
			writeUnauthorized(w)
			return
		}

		claims, err := security.ValidateAccessToken(m.jwtCfg, cookie.Value)
		if err != nil || m.repo == nil {
			writeUnauthorized(w)
			return
		}
		session, err := m.repo.GetSessionByJTI(r.Context(), claims.ID)
		if err != nil || !session.IsActive(time.Now()) {
			writeUnauthorized(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
}

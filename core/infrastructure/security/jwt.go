package security

import (
	"errors"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTConfig struct {
	Secret     string
	Issuer     string
	Audience   string
	Expiration time.Duration
	ClockSkew  time.Duration
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTConfig(secret, issuer, audience string, expiration time.Duration) JWTConfig {
	return JWTConfig{
		Secret:     secret,
		Issuer:     issuer,
		Audience:   audience,
		Expiration: expiration,
		ClockSkew:  30 * time.Second,
	}
}

func GenerateAccessToken(cfg JWTConfig, username string) (string, error) {
	if cfg.Secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	now := time.Now()
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    cfg.Issuer,
			Audience:  jwt.ClaimStrings{cfg.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.Expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ValidateAccessToken(cfg JWTConfig, tokenString string) (*Claims, error) {
	if cfg.Secret == "" {
		return nil, errors.New("jwt secret is empty")
	}

	parsed, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	if claims.ID == "" {
		return nil, errors.New("token id is missing")
	}

	// Manual validation for exp/nbf/iat (jwt.RegisteredClaims API differs across versions).
	now := time.Now()
	if claims.ExpiresAt != nil {
		if now.After(claims.ExpiresAt.Time.Add(cfg.ClockSkew)) {
			return nil, errors.New("token expired")
		}
	}
	if claims.NotBefore != nil {
		if now.Before(claims.NotBefore.Time.Add(-cfg.ClockSkew)) {
			return nil, errors.New("token not active")
		}
	}
	if claims.IssuedAt != nil {
		// no strict iat validation; kept for completeness
		_ = claims.IssuedAt
	}

	if cfg.Issuer != "" && claims.Issuer != cfg.Issuer {
		return nil, errors.New("invalid issuer")
	}
	if cfg.Audience != "" {
		found := slices.Contains(claims.Audience, cfg.Audience)
		if !found {
			return nil, errors.New("invalid audience")
		}
	}

	return claims, nil
}

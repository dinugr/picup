package authhandler

import (
	"testing"
	"time"

	"picup/core/infrastructure/security"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	cfg := security.NewJWTConfig("secret", "issuer", "aud", 10*time.Second)
	token, err := security.GenerateAccessToken(cfg, "admin")
	if err != nil {
		t.Fatalf("expected token generation success: %v", err)
	}

	claims, err := security.ValidateAccessToken(cfg, token)
	if err != nil {
		t.Fatalf("expected token validation success: %v", err)
	}
	if claims.Username != "admin" {
		t.Fatalf("expected username admin, got %s", claims.Username)
	}
	if claims.ID == "" {
		t.Fatal("expected a unique token id")
	}

	secondToken, err := security.GenerateAccessToken(cfg, "admin")
	if err != nil {
		t.Fatalf("expected second token generation success: %v", err)
	}
	secondClaims, err := security.ValidateAccessToken(cfg, secondToken)
	if err != nil {
		t.Fatalf("expected second token validation success: %v", err)
	}
	if claims.ID == secondClaims.ID {
		t.Fatal("expected token ids to be unique")
	}
}

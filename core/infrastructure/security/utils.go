package security

import (
	"crypto/rand"
	"encoding/base64"
)

func NewSecret(strlen int) string {
	secret := make([]byte, strlen)
	if _, err := rand.Read(secret); err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(secret)
}

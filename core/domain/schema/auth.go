package schema

import "time"

type Auth struct {
	JWTSecret            string
	JWTIssuer            string
	JWTAudience          string
	JWTExpirationSeconds int64
	JWTSecretAutogen     bool

	AdminUsername string
	AdminPassword string

	// reserved for future
	_ time.Duration
}

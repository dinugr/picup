package models

import "time"

type Session struct {
	ID               string
	UserID           string
	TokenJTI         string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	LastSeenAt       *time.Time
	UserAgent        string
	IPAddress        string
	RevokedAt        *time.Time
	RevocationReason string
}

func (s *Session) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && now.Before(s.ExpiresAt)
}

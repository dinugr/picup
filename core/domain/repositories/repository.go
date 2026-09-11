package repositories

import (
	"context"
	"picup/core/domain/models"
	"time"
)

type Repository interface {
	// images
	SaveImage(ctx context.Context, img *models.Image) error
	GetImageByID(ctx context.Context, id string) (*models.Image, error)
	ListImage(ctx context.Context, query models.ImageQuery) ([]*models.Image, int, error)
	ListImageByID(ctx context.Context, id string) ([]*models.Image, error)
	DeleteImage(ctx context.Context, id string) error

	// sessions
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByJTI(ctx context.Context, tokenJTI string) (*models.Session, error)
	ListActiveSessions(ctx context.Context, userID string, now time.Time) ([]*models.Session, error)
	RevokeSession(ctx context.Context, sessionID, userID, reason string, now time.Time) error
	RevokeSessionByJTI(ctx context.Context, tokenJTI, userID, reason string, now time.Time) error
	RevokeAllSessions(ctx context.Context, userID, reason string, now time.Time) error

	// misc
	Close() error
}

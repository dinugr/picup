package sqlite

import (
	"context"
	"testing"
	"time"

	"picup/core/domain/models"
)

func TestSessionLifecycleAndOwnership(t *testing.T) {
	repo, err := NewSQLiteRepository("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	now := time.Now().UTC().Truncate(time.Microsecond)
	ctx := context.Background()
	first := &models.Session{ID: "s1", UserID: "alice", TokenJTI: "jti-1", CreatedAt: now, ExpiresAt: now.Add(time.Hour), UserAgent: "browser", IPAddress: "127.0.0.1"}
	second := &models.Session{ID: "s2", UserID: "alice", TokenJTI: "jti-2", CreatedAt: now, ExpiresAt: now.Add(time.Hour)}
	otherUser := &models.Session{ID: "s3", UserID: "bob", TokenJTI: "jti-3", CreatedAt: now, ExpiresAt: now.Add(time.Hour)}
	for _, session := range []*models.Session{first, second, otherUser} {
		if err := repo.CreateSession(ctx, session); err != nil {
			t.Fatal(err)
		}
	}

	found, err := repo.GetSessionByJTI(ctx, "jti-1")
	if err != nil || found.UserID != "alice" {
		t.Fatalf("lookup failed: session=%+v err=%v", found, err)
	}

	if err := repo.RevokeSession(ctx, "s1", "bob", "logout", now); err != nil {
		t.Fatal(err)
	}
	active, err := repo.ListActiveSessions(ctx, "alice", now)
	if err != nil || len(active) != 2 {
		t.Fatalf("wrong active sessions after unauthorized revoke: got %d err=%v", len(active), err)
	}

	if err := repo.RevokeSession(ctx, "s1", "alice", "logout", now); err != nil {
		t.Fatal(err)
	}
	active, err = repo.ListActiveSessions(ctx, "alice", now)
	if err != nil || len(active) != 1 || active[0].ID != "s2" {
		t.Fatalf("wrong active sessions after revoke: got %+v err=%v", active, err)
	}

	if err := repo.RevokeAllSessions(ctx, "alice", "logout_all", now); err != nil {
		t.Fatal(err)
	}
	active, err = repo.ListActiveSessions(ctx, "alice", now)
	if err != nil || len(active) != 0 {
		t.Fatalf("expected alice sessions revoked: got %+v err=%v", active, err)
	}
	active, err = repo.ListActiveSessions(ctx, "bob", now)
	if err != nil || len(active) != 1 || active[0].ID != "s3" {
		t.Fatalf("logout-all affected another user: got %+v err=%v", active, err)
	}
}

func TestListActiveSessionsExcludesExpired(t *testing.T) {
	repo, err := NewSQLiteRepository("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	now := time.Now().UTC()
	if err := repo.CreateSession(context.Background(), &models.Session{ID: "expired", UserID: "alice", TokenJTI: "expired-jti", CreatedAt: now.Add(-2 * time.Hour), ExpiresAt: now.Add(-time.Minute)}); err != nil {
		t.Fatal(err)
	}
	active, err := repo.ListActiveSessions(context.Background(), "alice", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(active) != 0 {
		t.Fatalf("expected expired session excluded, got %d", len(active))
	}
}

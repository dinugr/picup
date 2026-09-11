package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"picup/core/domain/models"
	"strings"
	"time"
)

func (r *SQLiteRepository) SaveImage(ctx context.Context, img *models.Image) error {
	query := `
	INSERT INTO images (id, created_at, upload_name, stored_name, type, variant, master_id, mime_type, size_bytes, width, height)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		img.ID,
		img.CreatedAt.Format(time.RFC3339),
		img.UploadName,
		img.StoredName,
		img.Type,
		img.Variant,
		img.MasterID,
		img.MIMEType,
		img.SizeBytes,
		img.Width,
		img.Height,
	)
	if err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) GetImageByID(ctx context.Context, id string) (*models.Image, error) {
	query := `
	SELECT id, created_at, upload_name, stored_name, type, variant, master_id, mime_type, size_bytes, width, height
	FROM images
	WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var img models.Image
	var createdAtStr string
	err := row.Scan(
		&img.ID,
		&createdAtStr,
		&img.UploadName,
		&img.StoredName,
		&img.Type,
		&img.Variant,
		&img.MasterID,
		&img.MIMEType,
		&img.SizeBytes,
		&img.Width,
		&img.Height,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image not found: %s", id)
		}
		return nil, fmt.Errorf("failed to query image: %w", err)
	}

	t, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		// Fallback parse
		t, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}
	img.CreatedAt = t

	return &img, nil
}

func (r *SQLiteRepository) ListImage(ctx context.Context, query models.ImageQuery) ([]*models.Image, int, error) {
	whereClause := " WHERE 1=1"
	var args []any

	// 1. ID Filter
	if query.ID != "" {
		whereClause += " AND (id = ? OR master_id = ?)"
		args = append(args, query.ID, query.ID)
	}

	// 2. Type Filter
	if len(query.Type) > 0 {
		placeholders := make([]string, len(query.Type))
		for i, v := range query.Type {
			placeholders[i] = "?"
			args = append(args, v)
		}

		inverseFlag := ""
		if query.TypeExclude {
			inverseFlag = " NOT"
		}

		whereClause += fmt.Sprintf(" AND type%s IN (%s)", inverseFlag, strings.Join(placeholders, ","))
	}

	// 3. Search Filter
	if query.Search != "" {
		if strings.Contains(query.Search, "/") {
			// Exact match for MIME types
			whereClause += " AND mime_type = ?"
			args = append(args, query.Search)
		} else {
			// Partial match for filenames
			whereClause += " AND (upload_name LIKE ? OR stored_name LIKE ?)"
			searchTerm := "%" + query.Search + "%"
			args = append(args, searchTerm, searchTerm)
		}
	}

	// 4. Get Total Count (Executes BEFORE pagination/sorting)
	countQuery := "SELECT COUNT(*) FROM images" + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count images: %w", err)
	}

	// 5. Sort By & Sort Order
	sortBy := "created_at"
	validSortColumns := map[string]bool{
		"created_at":  true,
		"upload_name": true,
		"size_bytes":  true,
		"type":        true,
	}
	if query.SortBy != "" && validSortColumns[query.SortBy] {
		sortBy = query.SortBy
	}

	sortOrder := "ASC"
	if strings.ToUpper(query.SortOrder) == "DESC" {
		sortOrder = "DESC"
	}

	orderByClause := fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// 6. Pagination
	limitClause := ""
	if query.PageIndex >= 0 {
		offset := math.Max(float64((query.PageIndex-1)*query.PageSize), 0)

		limitClause = " LIMIT ? OFFSET ?"
		args = append(args, query.PageSize, offset) // Safe to append here since Count query is done
	}

	// 7. Execute Main Query
	selectClause := "SELECT id, created_at, upload_name, stored_name, type, variant, master_id, mime_type, size_bytes, width, height FROM images"
	finalQuery := selectClause + whereClause + orderByClause + limitClause

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get images: %w", err)
	}
	defer rows.Close()

	var images []*models.Image
	for rows.Next() {
		var img models.Image
		var createdAtStr string
		if err := rows.Scan(
			&img.ID,
			&createdAtStr,
			&img.UploadName,
			&img.StoredName,
			&img.Type,
			&img.Variant,
			&img.MasterID,
			&img.MIMEType,
			&img.SizeBytes,
			&img.Width,
			&img.Height,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan image row: %w", err)
		}

		t, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			t, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		img.CreatedAt = t

		images = append(images, &img)
	}

	return images, total, nil
}

func (r *SQLiteRepository) ListImageBySourceID(ctx context.Context, sourceID string) ([]*models.Image, error) {
	query := `
	SELECT id, created_at, upload_name, stored_name, type, variant, master_id, mime_type, size_bytes, width, height
	FROM images
	WHERE master_id = ?
	ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list variants by source id: %w", err)
	}
	defer rows.Close()

	var images []*models.Image
	for rows.Next() {
		var img models.Image
		var createdAtStr string
		if err := rows.Scan(
			&img.ID,
			&createdAtStr,
			&img.UploadName,
			&img.StoredName,
			&img.Type,
			&img.Variant,
			&img.MasterID,
			&img.MIMEType,
			&img.SizeBytes,
			&img.Width,
			&img.Height,
		); err != nil {
			return nil, fmt.Errorf("failed to scan image row: %w", err)
		}

		t, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			t, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		img.CreatedAt = t

		images = append(images, &img)
	}

	return images, nil
}

func (r *SQLiteRepository) ListImageByID(ctx context.Context, ID string) ([]*models.Image, error) {
	query := `
	SELECT id, created_at, upload_name, stored_name, type, variant, master_id, mime_type, size_bytes, width, height
	FROM images
	WHERE master_id = ? or id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, ID, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list variants by source id: %w", err)
	}
	defer rows.Close()

	var images []*models.Image
	for rows.Next() {
		var img models.Image
		var createdAtStr string
		if err := rows.Scan(
			&img.ID,
			&createdAtStr,
			&img.UploadName,
			&img.StoredName,
			&img.Type,
			&img.Variant,
			&img.MasterID,
			&img.MIMEType,
			&img.SizeBytes,
			&img.Width,
			&img.Height,
		); err != nil {
			return nil, fmt.Errorf("failed to scan image row: %w", err)
		}

		t, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			t, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		img.CreatedAt = t

		images = append(images, &img)
	}

	return images, nil
}

func (r *SQLiteRepository) DeleteImage(ctx context.Context, id string) error {
	query := `DELETE FROM images WHERE id = ? or master_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, id)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("image not found: %s", id)
	}
	return nil
}

func (r *SQLiteRepository) CreateSession(ctx context.Context, session *models.Session) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth_sessions (id, user_id, token_jti, created_at, expires_at, last_seen_at, user_agent, ip_address, revoked_at, revocation_reason)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, session.ID, session.UserID, session.TokenJTI, session.CreatedAt.UTC(), session.ExpiresAt.UTC(), nullableTime(session.LastSeenAt), session.UserAgent, session.IPAddress, nullableTime(session.RevokedAt), session.RevocationReason)
	if err != nil {
		return fmt.Errorf("failed to create auth session: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) GetSessionByJTI(ctx context.Context, tokenJTI string) (*models.Session, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_jti, created_at, expires_at, last_seen_at, user_agent, ip_address, revoked_at, revocation_reason
		FROM auth_sessions WHERE token_jti = ?
	`, tokenJTI)
	session, err := scanSession(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("auth session not found: %w", err)
		}
		return nil, fmt.Errorf("failed to query auth session: %w", err)
	}
	return session, nil
}

func (r *SQLiteRepository) ListActiveSessions(ctx context.Context, userID string, now time.Time) ([]*models.Session, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, token_jti, created_at, expires_at, last_seen_at, user_agent, ip_address, revoked_at, revocation_reason
		FROM auth_sessions
		WHERE user_id = ? AND revoked_at IS NULL AND expires_at > ?
		ORDER BY created_at DESC
	`, userID, now.UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to list auth sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session, err := scanSession(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan auth session: %w", err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate auth sessions: %w", err)
	}
	return sessions, nil
}

func (r *SQLiteRepository) RevokeSession(ctx context.Context, sessionID, userID, reason string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE auth_sessions SET revoked_at = ?, revocation_reason = ?
		WHERE id = ? AND user_id = ? AND revoked_at IS NULL
	`, now.UTC(), reason, sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke auth session: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) RevokeSessionByJTI(ctx context.Context, tokenJTI, userID, reason string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE auth_sessions SET revoked_at = ?, revocation_reason = ?
		WHERE token_jti = ? AND user_id = ? AND revoked_at IS NULL
	`, now.UTC(), reason, tokenJTI, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke auth session: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) RevokeAllSessions(ctx context.Context, userID, reason string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE auth_sessions SET revoked_at = ?, revocation_reason = ?
		WHERE user_id = ? AND revoked_at IS NULL
	`, now.UTC(), reason, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all auth sessions: %w", err)
	}
	return nil
}

type sessionScanner interface {
	Scan(dest ...any) error
}

func scanSession(scanner sessionScanner) (*models.Session, error) {
	var session models.Session
	var createdAt, expiresAt time.Time
	var lastSeenAt, revokedAt sql.NullTime
	if err := scanner.Scan(&session.ID, &session.UserID, &session.TokenJTI, &createdAt, &expiresAt, &lastSeenAt, &session.UserAgent, &session.IPAddress, &revokedAt, &session.RevocationReason); err != nil {
		return nil, err
	}
	session.CreatedAt = createdAt
	session.ExpiresAt = expiresAt
	if lastSeenAt.Valid {
		session.LastSeenAt = &lastSeenAt.Time
	}
	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.Time
	}
	return &session, nil
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.UTC()
}

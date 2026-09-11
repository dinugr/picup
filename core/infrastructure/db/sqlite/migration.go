package sqlite

func (r *SQLiteRepository) autoMigrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS images (
		id             TEXT PRIMARY KEY,
		created_at     DATETIME NOT NULL,
		upload_name    TEXT NOT NULL,
		stored_name    TEXT NOT NULL,
		type           TEXT NOT NULL,
		variant        TEXT NOT NULL,
		master_id      TEXT,
		mime_type      TEXT NOT NULL,
		size_bytes     INTEGER NOT NULL,
		width          INTEGER DEFAULT 0,
		height         INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_images_master_id ON images(master_id);
	CREATE INDEX IF NOT EXISTS idx_images_upload_name ON images(upload_name);

	CREATE TABLE IF NOT EXISTS auth_sessions (
		id                TEXT PRIMARY KEY,
		user_id           TEXT NOT NULL,
		token_jti         TEXT NOT NULL UNIQUE,
		created_at        DATETIME NOT NULL,
		expires_at        DATETIME NOT NULL,
		last_seen_at      DATETIME,
		user_agent        TEXT NOT NULL DEFAULT '',
		ip_address        TEXT NOT NULL DEFAULT '',
		revoked_at        DATETIME,
		revocation_reason TEXT NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_active ON auth_sessions(user_id, revoked_at, expires_at);
	CREATE INDEX IF NOT EXISTS idx_auth_sessions_expires_at ON auth_sessions(expires_at);
	`
	_, err := r.db.Exec(query)
	return err
}

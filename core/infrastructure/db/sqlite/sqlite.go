package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(params string) (*SQLiteRepository, error) {

	source := params
	if source == "" {
		return nil, fmt.Errorf("require \"source\" field in db configuration")
	}

	db, err := sql.Open("sqlite3", source)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database (%s): %w", source, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.autoMigrate(); err != nil {
		return nil, fmt.Errorf("sqlite auto migration failed: %w", err)
	}

	return repo, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

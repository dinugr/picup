package db

import (
	"fmt"
	"picup/core/domain/repositories"
	"picup/core/infrastructure/db/sqlite"
	"strings"
)

func NewRepository(dsn string) (repositories.Repository, error) {

	driver, args := parseDsnString(dsn)
	switch driver {
	case "sqlite", "sqlite3":
		return sqlite.NewSQLiteRepository(args)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

// parseDsnString splits a DSN string formatted as driver_name://arguments...
// into its driver name and argument components.
func parseDsnString(dsn string) (string, string) {
	driver, args, found := strings.Cut(dsn, "://")
	if !found {
		// Return empty strings or handle unsupported format
		return "", dsn
	}
	return driver, args
}

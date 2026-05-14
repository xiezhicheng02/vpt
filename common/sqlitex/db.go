package sqlitex

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// Open creates the directory of dbPath (if missing) and opens a sqlite database.
// Caller imports the desired sqlite driver, e.g. _ "modernc.org/sqlite".
func Open(driver, dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB, ddls []string) error {
	for _, ddl := range ddls {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("migrate: %w (sql=%s)", err, ddl)
		}
	}
	return nil
}

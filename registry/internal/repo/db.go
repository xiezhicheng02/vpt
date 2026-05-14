package repo

import (
	"database/sql"

	"github.com/vpt/common/sqlitex"
	_ "modernc.org/sqlite"
)

var ddls = []string{
	`CREATE TABLE IF NOT EXISTS services (
		instance_id TEXT PRIMARY KEY,
		name        TEXT NOT NULL,
		address     TEXT NOT NULL,
		scheme      TEXT NOT NULL,
		status      TEXT NOT NULL,
		last_seen   DATETIME NOT NULL
	)`,
	`CREATE INDEX IF NOT EXISTS idx_services_name ON services(name)`,
	`CREATE TABLE IF NOT EXISTS configs (
		key        TEXT PRIMARY KEY,
		value      TEXT NOT NULL,
		updated_at DATETIME NOT NULL
	)`,
}

func OpenDB(path string) (*sql.DB, error) {
	db, err := sqlitex.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := sqlitex.Migrate(db, ddls); err != nil {
		return nil, err
	}
	return db, nil
}

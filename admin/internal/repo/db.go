package repo

import (
	"database/sql"

	"github.com/vpt/common/sqlitex"
	_ "modernc.org/sqlite"
)

var ddls = []string{
	`CREATE TABLE IF NOT EXISTS users (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		email          TEXT UNIQUE NOT NULL,
		username       TEXT UNIQUE NOT NULL,
		password_hash  TEXT NOT NULL,
		passkey        TEXT UNIQUE NOT NULL,
		role           TEXT NOT NULL DEFAULT 'user',
		email_verified INTEGER NOT NULL DEFAULT 0,
		created_at     DATETIME NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS tokens (
		token      TEXT PRIMARY KEY,
		user_id    INTEGER NOT NULL,
		kind       TEXT NOT NULL,
		expires_at DATETIME NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS torrents (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		info_hash   TEXT UNIQUE NOT NULL,
		name        TEXT NOT NULL,
		size        INTEGER NOT NULL DEFAULT 0,
		uploader_id INTEGER NOT NULL,
		file_path   TEXT NOT NULL,
		created_at  DATETIME NOT NULL
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

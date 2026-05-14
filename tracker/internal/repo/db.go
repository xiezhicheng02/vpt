package repo

import (
	"database/sql"

	"github.com/vpt/common/sqlitex"
	_ "modernc.org/sqlite"
)

var ddls = []string{
	`CREATE TABLE IF NOT EXISTS peers (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		info_hash   TEXT NOT NULL,
		peer_id     TEXT NOT NULL,
		user_id     INTEGER NOT NULL,
		ip          TEXT NOT NULL,
		port        INTEGER NOT NULL,
		uploaded    INTEGER NOT NULL DEFAULT 0,
		downloaded  INTEGER NOT NULL DEFAULT 0,
		left_bytes  INTEGER NOT NULL DEFAULT 0,
		event       TEXT NOT NULL DEFAULT '',
		last_seen   DATETIME NOT NULL,
		UNIQUE(info_hash, peer_id)
	)`,
	`CREATE INDEX IF NOT EXISTS idx_peers_info_hash ON peers(info_hash)`,
	`CREATE TABLE IF NOT EXISTS stats (
		info_hash   TEXT PRIMARY KEY,
		seeders     INTEGER NOT NULL DEFAULT 0,
		leechers    INTEGER NOT NULL DEFAULT 0,
		completed   INTEGER NOT NULL DEFAULT 0,
		updated_at  DATETIME NOT NULL
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

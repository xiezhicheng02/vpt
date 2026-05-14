package repo

import (
	"database/sql"
	"time"

	"github.com/vpt/tracker/internal/model"
)

type StatRepo struct{ db *sql.DB }

func NewStatRepo(db *sql.DB) *StatRepo { return &StatRepo{db: db} }

func (r *StatRepo) Update(infoHash string, seeders, leechers int, completedDelta int64) error {
	_, err := r.db.Exec(`INSERT INTO stats(info_hash,seeders,leechers,completed,updated_at)
		VALUES(?,?,?,?,?)
		ON CONFLICT(info_hash) DO UPDATE SET
		  seeders=excluded.seeders, leechers=excluded.leechers,
		  completed=stats.completed+?, updated_at=excluded.updated_at`,
		infoHash, seeders, leechers, completedDelta, time.Now(), completedDelta)
	return err
}

func (r *StatRepo) Get(infoHash string) (*model.StatSnapshot, error) {
	row := r.db.QueryRow(`SELECT info_hash,seeders,leechers,completed,updated_at FROM stats WHERE info_hash=?`, infoHash)
	var s model.StatSnapshot
	if err := row.Scan(&s.InfoHash, &s.Seeders, &s.Leechers, &s.Completed, &s.UpdatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StatRepo) List(limit, offset int) ([]model.StatSnapshot, error) {
	rows, err := r.db.Query(`SELECT info_hash,seeders,leechers,completed,updated_at FROM stats ORDER BY updated_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.StatSnapshot
	for rows.Next() {
		var s model.StatSnapshot
		if err := rows.Scan(&s.InfoHash, &s.Seeders, &s.Leechers, &s.Completed, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

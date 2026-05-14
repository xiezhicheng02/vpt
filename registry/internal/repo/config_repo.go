package repo

import (
	"database/sql"
	"time"

	"github.com/vpt/registry/internal/model"
)

type ConfigRepo struct{ db *sql.DB }

func NewConfigRepo(db *sql.DB) *ConfigRepo { return &ConfigRepo{db: db} }

func (r *ConfigRepo) Get(key string) (*model.ConfigItem, error) {
	row := r.db.QueryRow(`SELECT key,value,updated_at FROM configs WHERE key=?`, key)
	var c model.ConfigItem
	if err := row.Scan(&c.Key, &c.Value, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConfigRepo) Set(key, value string) error {
	_, err := r.db.Exec(`INSERT INTO configs(key,value,updated_at) VALUES(?,?,?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value, updated_at=excluded.updated_at`,
		key, value, time.Now())
	return err
}

func (r *ConfigRepo) List() ([]model.ConfigItem, error) {
	rows, err := r.db.Query(`SELECT key,value,updated_at FROM configs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.ConfigItem
	for rows.Next() {
		var c model.ConfigItem
		if err := rows.Scan(&c.Key, &c.Value, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

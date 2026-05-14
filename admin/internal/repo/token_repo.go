package repo

import (
	"database/sql"

	"github.com/vpt/admin/internal/model"
)

type TokenRepo struct{ db *sql.DB }

func NewTokenRepo(db *sql.DB) *TokenRepo { return &TokenRepo{db: db} }

func (r *TokenRepo) Create(t *model.Token) error {
	_, err := r.db.Exec(`INSERT INTO tokens(token,user_id,kind,expires_at) VALUES(?,?,?,?)`,
		t.Token, t.UserID, t.Kind, t.ExpiresAt)
	return err
}

func (r *TokenRepo) Find(token string) (*model.Token, error) {
	row := r.db.QueryRow(`SELECT token,user_id,kind,expires_at FROM tokens WHERE token=?`, token)
	var t model.Token
	if err := row.Scan(&t.Token, &t.UserID, &t.Kind, &t.ExpiresAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TokenRepo) Delete(token string) error {
	_, err := r.db.Exec(`DELETE FROM tokens WHERE token=?`, token)
	return err
}

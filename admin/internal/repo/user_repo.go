package repo

import (
	"database/sql"

	"github.com/vpt/admin/internal/model"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(u *model.User) error {
	res, err := r.db.Exec(`INSERT INTO users(email,username,password_hash,passkey,role,email_verified,created_at)
		VALUES(?,?,?,?,?,?,?)`,
		u.Email, u.Username, u.PasswordHash, u.Passkey, u.Role, u.EmailVerified, u.CreatedAt)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	u.ID = id
	return nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	return r.scanOne(`SELECT id,email,username,password_hash,passkey,role,email_verified,created_at FROM users WHERE email=?`, email)
}

func (r *UserRepo) FindByID(id int64) (*model.User, error) {
	return r.scanOne(`SELECT id,email,username,password_hash,passkey,role,email_verified,created_at FROM users WHERE id=?`, id)
}

func (r *UserRepo) FindByPasskey(passkey string) (*model.User, error) {
	return r.scanOne(`SELECT id,email,username,password_hash,passkey,role,email_verified,created_at FROM users WHERE passkey=?`, passkey)
}

func (r *UserRepo) MarkVerified(id int64) error {
	_, err := r.db.Exec(`UPDATE users SET email_verified=1 WHERE id=?`, id)
	return err
}

func (r *UserRepo) scanOne(query string, args ...any) (*model.User, error) {
	row := r.db.QueryRow(query, args...)
	var u model.User
	var verified int
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Passkey, &u.Role, &verified, &u.CreatedAt); err != nil {
		return nil, err
	}
	u.EmailVerified = verified == 1
	return &u, nil
}

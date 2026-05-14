package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Passkey      string    `json:"passkey"` // used by tracker to identify user
	Role         string    `json:"role"`    // user / admin
	EmailVerified bool     `json:"email_verified"`
	CreatedAt    time.Time `json:"created_at"`
}

type Token struct {
	Token     string    `json:"token"`
	UserID    int64     `json:"user_id"`
	Kind      string    `json:"kind"` // session / email_verify
	ExpiresAt time.Time `json:"expires_at"`
}

type Torrent struct {
	ID         int64     `json:"id"`
	InfoHash   string    `json:"info_hash"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	UploaderID int64     `json:"uploader_id"`
	FilePath   string    `json:"file_path"`
	CreatedAt  time.Time `json:"created_at"`
}

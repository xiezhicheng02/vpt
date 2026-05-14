package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/vpt/admin/internal/config"
	"github.com/vpt/admin/internal/model"
	"github.com/vpt/admin/internal/repo"
)

type User struct {
	users  *repo.UserRepo
	tokens *repo.TokenRepo
	cfg    *config.Config
}

func NewUser(u *repo.UserRepo, t *repo.TokenRepo, cfg *config.Config) *User {
	return &User{users: u, tokens: t, cfg: cfg}
}

func (s *User) Register(email, username, password string) (*model.User, string, error) {
	if email == "" || username == "" || password == "" {
		return nil, "", errors.New("missing fields")
	}
	u := &model.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashPassword(password, s.cfg.TokenSecret),
		Passkey:      randomHex(20),
		Role:         "user",
		CreatedAt:    time.Now(),
	}
	if err := s.users.Create(u); err != nil {
		return nil, "", err
	}
	verifyToken := randomHex(24)
	_ = s.tokens.Create(&model.Token{
		Token:     verifyToken,
		UserID:    u.ID,
		Kind:      "email_verify",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	// TODO: send verification email via SMTP using cfg.SMTP*
	return u, verifyToken, nil
}

func (s *User) Login(email, password string) (string, error) {
	u, err := s.users.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if u.PasswordHash != hashPassword(password, s.cfg.TokenSecret) {
		return "", errors.New("invalid credentials")
	}
	token := randomHex(32)
	err = s.tokens.Create(&model.Token{
		Token:     token,
		UserID:    u.ID,
		Kind:      "session",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	return token, err
}

func (s *User) Verify(token string) (*model.User, error) {
	t, err := s.tokens.Find(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if t.Kind != "session" || time.Now().After(t.ExpiresAt) {
		return nil, errors.New("expired token")
	}
	return s.users.FindByID(t.UserID)
}

func (s *User) VerifyEmail(token string) error {
	t, err := s.tokens.Find(token)
	if err != nil {
		return errors.New("invalid token")
	}
	if t.Kind != "email_verify" || time.Now().After(t.ExpiresAt) {
		return errors.New("expired token")
	}
	if err := s.users.MarkVerified(t.UserID); err != nil {
		return err
	}
	_ = s.tokens.Delete(token)
	return nil
}

func (s *User) FindByPasskey(passkey string) (*model.User, error) {
	return s.users.FindByPasskey(passkey)
}

func hashPassword(pw, secret string) string {
	h := sha256.Sum256([]byte(pw + ":" + secret))
	return hex.EncodeToString(h[:])
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

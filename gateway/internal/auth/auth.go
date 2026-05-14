package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/vpt/common/httpx"
	"github.com/vpt/common/registry"
)

type Client struct {
	reg *registry.Client
}

func NewClient(reg *registry.Client) *Client { return &Client{reg: reg} }

type VerifyResult struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Verify calls admin service /api/v1/auth/verify with the bearer token.
func (c *Client) Verify(ctx context.Context, token string) (*VerifyResult, error) {
	instances, err := c.reg.Discover(ctx, "admin")
	if err != nil || len(instances) == 0 {
		return nil, errors.New("admin service unavailable")
	}
	base := instances[0].Scheme + "://" + instances[0].Address
	hc := httpx.NewClient(base)
	var resp VerifyResult
	if err := hc.PostJSON(ctx, "/api/v1/auth/verify", map[string]string{"token": token}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Middleware enforces auth on protected routes.
func (c *Client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		res, err := c.Verify(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-User-Id", itoa(res.UserID))
		r.Header.Set("X-Username", res.Username)
		r.Header.Set("X-Role", res.Role)
		next.ServeHTTP(w, r)
	})
}

func itoa(i int64) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	n := len(b)
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		n--
		b[n] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		n--
		b[n] = '-'
	}
	return string(b[n:])
}

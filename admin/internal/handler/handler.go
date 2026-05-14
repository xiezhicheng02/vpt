package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vpt/common/httpx"
	"github.com/vpt/admin/internal/service"
)

func RegisterRoutes(mux *http.ServeMux, userSvc *service.User, torrentSvc *service.Torrent) {
	mux.HandleFunc("/api/v1/auth/register", register(userSvc))
	mux.HandleFunc("/api/v1/auth/login", login(userSvc))
	mux.HandleFunc("/api/v1/auth/verify", verify(userSvc))         // gateway internal RPC
	mux.HandleFunc("/api/v1/auth/verify-email", verifyEmail(userSvc))
	mux.HandleFunc("/api/v1/auth/passkey", passkeyLookup(userSvc)) // tracker internal RPC

	mux.HandleFunc("/api/v1/torrents/upload", uploadTorrent(torrentSvc))
	mux.HandleFunc("/api/v1/torrents/list", listTorrents(torrentSvc))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { httpx.OK(w, "ok") })
}

func register(s *service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Username, Password string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Fail(w, 400, 1001, "bad request")
			return
		}
		u, vt, err := s.Register(req.Email, req.Username, req.Password)
		if err != nil {
			httpx.Fail(w, 400, 1002, err.Error())
			return
		}
		httpx.OK(w, map[string]any{"user_id": u.ID, "verify_token": vt})
	}
}

func login(s *service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Password string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Fail(w, 400, 1001, "bad request")
			return
		}
		token, err := s.Login(req.Email, req.Password)
		if err != nil {
			httpx.Fail(w, 401, 4011, err.Error())
			return
		}
		httpx.OK(w, map[string]string{"token": token})
	}
}

func verify(s *service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Fail(w, 400, 1001, "bad request")
			return
		}
		u, err := s.Verify(req.Token)
		if err != nil {
			httpx.Fail(w, 401, 4011, err.Error())
			return
		}
		// Direct JSON shape (gateway expects user_id/username/role)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"user_id":  u.ID,
			"username": u.Username,
			"role":     u.Role,
		})
	}
}

func verifyEmail(s *service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if err := s.VerifyEmail(token); err != nil {
			httpx.Fail(w, 400, 1003, err.Error())
			return
		}
		httpx.OK(w, nil)
	}
}

func passkeyLookup(s *service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pk := r.URL.Query().Get("passkey")
		u, err := s.FindByPasskey(pk)
		if err != nil {
			httpx.Fail(w, 404, 4041, "not found")
			return
		}
		httpx.OK(w, map[string]any{"user_id": u.ID, "username": u.Username})
	}
}

func uploadTorrent(s *service.Torrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			httpx.Fail(w, 400, 1001, err.Error())
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.Fail(w, 400, 1001, "missing file")
			return
		}
		defer file.Close()
		uid, _ := strconv.ParseInt(r.Header.Get("X-User-Id"), 10, 64)
		t, err := s.Upload(uid, header.Filename, file)
		if err != nil {
			httpx.Fail(w, 500, 5001, err.Error())
			return
		}
		httpx.OK(w, t)
	}
}

func listTorrents(s *service.Torrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 {
			limit = 20
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		list, err := s.List(limit, offset)
		if err != nil {
			httpx.Fail(w, 500, 5001, err.Error())
			return
		}
		httpx.OK(w, list)
	}
}

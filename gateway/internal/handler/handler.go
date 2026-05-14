package handler

import (
	"net/http"

	"github.com/vpt/common/httpx"
	"github.com/vpt/gateway/internal/auth"
	"github.com/vpt/gateway/internal/proxy"
)

func RegisterRoutes(mux *http.ServeMux, router *proxy.Router, authClient *auth.Client) {
	// Public: registration / login endpoints on admin service
	mux.Handle("/api/v1/auth/register", router.ProxyHTTP("admin", ""))
	mux.Handle("/api/v1/auth/login", router.ProxyHTTP("admin", ""))
	mux.Handle("/api/v1/auth/verify-email", router.ProxyHTTP("admin", ""))

	// Public: tracker announce/scrape (BT clients send these directly)
	mux.Handle("/announce", router.ProxyHTTP("tracker", ""))
	mux.Handle("/scrape", router.ProxyHTTP("tracker", ""))

	// Protected: admin (torrents/users) — require auth
	mux.Handle("/api/v1/torrents/", authClient.Middleware(router.ProxyHTTP("admin", "")))
	mux.Handle("/api/v1/users/", authClient.Middleware(router.ProxyHTTP("admin", "")))

	// Protected: tracker stats
	mux.Handle("/api/v1/stats/", authClient.Middleware(router.ProxyHTTP("tracker", "")))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { httpx.OK(w, "ok") })
}

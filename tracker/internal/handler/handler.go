package handler

import (
	"net"
	"net/http"
	"strconv"

	"github.com/vpt/common/httpx"
	"github.com/vpt/tracker/internal/model"
	"github.com/vpt/tracker/internal/service"
)

func RegisterRoutes(mux *http.ServeMux, t *service.Tracker) {
	mux.HandleFunc("/announce", announce(t))
	mux.HandleFunc("/scrape", scrape(t))
	mux.HandleFunc("/api/v1/stats/list", listStats(t))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { httpx.OK(w, "ok") })
}

// announce handles BT HTTP tracker protocol (BEP 3).
// TODO: encode response as bencoded dict; current placeholder returns JSON for testing.
func announce(t *service.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		port, _ := strconv.Atoi(q.Get("port"))
		up, _ := strconv.ParseInt(q.Get("uploaded"), 10, 64)
		down, _ := strconv.ParseInt(q.Get("downloaded"), 10, 64)
		left, _ := strconv.ParseInt(q.Get("left"), 10, 64)
		numWant, _ := strconv.Atoi(q.Get("numwant"))

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		req := &model.AnnounceRequest{
			InfoHash:   q.Get("info_hash"),
			PeerID:     q.Get("peer_id"),
			Passkey:    q.Get("passkey"),
			IP:         ip,
			Port:       port,
			Uploaded:   up,
			Downloaded: down,
			Left:       left,
			Event:      q.Get("event"),
			NumWant:    numWant,
			Compact:    q.Get("compact") == "1",
		}
		resp, err := t.Announce(r.Context(), req)
		if err != nil {
			httpx.Fail(w, 400, 1001, err.Error())
			return
		}
		httpx.OK(w, resp)
	}
}

func scrape(t *service.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Query().Get("info_hash")
		s, err := t.Scrape(hash)
		if err != nil {
			httpx.Fail(w, 404, 4041, "not found")
			return
		}
		httpx.OK(w, s)
	}
}

func listStats(t *service.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 {
			limit = 50
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		list, err := t.ListStats(limit, offset)
		if err != nil {
			httpx.Fail(w, 500, 5001, err.Error())
			return
		}
		httpx.OK(w, list)
	}
}

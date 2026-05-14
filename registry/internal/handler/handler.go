package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vpt/common/httpx"
	"github.com/vpt/registry/internal/service"
)

func RegisterRoutes(mux *http.ServeMux, reg *service.Registry, cfg *service.Config) {
	mux.HandleFunc("/api/v1/registry/register", registerHandler(reg))
	mux.HandleFunc("/api/v1/registry/heartbeat", heartbeatHandler(reg))
	mux.HandleFunc("/api/v1/registry/discover", discoverHandler(reg))
	mux.HandleFunc("/api/v1/config/get", configGetHandler(cfg))
	mux.HandleFunc("/api/v1/config/set", configSetHandler(cfg))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { httpx.OK(w, "ok") })
}

func registerHandler(reg *service.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Name, Address, Scheme string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Fail(w, 400, 1001, "bad request")
			return
		}
		inst, err := reg.Register(req.Name, req.Address, req.Scheme)
		if err != nil {
			httpx.Fail(w, 500, 5001, err.Error())
			return
		}
		httpx.OK(w, map[string]string{"instance_id": inst.InstanceID})
	}
}

func heartbeatHandler(reg *service.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			InstanceID string `json:"instance_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Fail(w, 400, 1001, "bad request")
			return
		}
		if err := reg.Heartbeat(req.InstanceID); err != nil {
			httpx.Fail(w, 500, 5001, err.Error())
			return
		}
		httpx.OK(w, nil)
	}
}

func discoverHandler(reg *service.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		list, err := reg.Discover(name)
		if err != nil {
			httpx.Fail(w, 500, 5001, err.Error())
			return
		}
		httpx.OK(w, list)
	}
}

func configGetHandler(cfg *service.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		item, err := cfg.Get(key)
		if err != nil {
			httpx.Fail(w, 404, 4041, "not found")
			return
		}
		httpx.OK(w, map[string]string{"value": item.Value})
	}
}

func configSetHandler(cfg *service.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Key, Value string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Fail(w, 400, 1001, "bad request")
			return
		}
		if err := cfg.Set(req.Key, req.Value); err != nil {
			httpx.Fail(w, 500, 5001, err.Error())
			return
		}
		httpx.OK(w, nil)
	}
}

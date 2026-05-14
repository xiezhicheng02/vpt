package proxy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/vpt/common/registry"
)

// Router resolves a service name to an upstream URL via the registry and proxies HTTP requests.
type Router struct {
	reg     *registry.Client
	counter atomic.Uint64
}

func NewRouter(reg *registry.Client) *Router { return &Router{reg: reg} }

func (r *Router) pickUpstream(ctx context.Context, name string) (*url.URL, error) {
	instances, err := r.reg.Discover(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return nil, errors.New("no upstream: " + name)
	}
	idx := r.counter.Add(1) % uint64(len(instances))
	i := instances[idx]
	return url.Parse(i.Scheme + "://" + i.Address)
}

// ProxyHTTP forwards request to the named service, stripping pathPrefix.
func (r *Router) ProxyHTTP(service, pathPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		upstream, err := r.pickUpstream(req.Context(), service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		rp := httputil.NewSingleHostReverseProxy(upstream)
		origDirector := rp.Director
		rp.Director = func(r *http.Request) {
			origDirector(r)
			if pathPrefix != "" && len(r.URL.Path) >= len(pathPrefix) && r.URL.Path[:len(pathPrefix)] == pathPrefix {
				r.URL.Path = r.URL.Path[len(pathPrefix):]
			}
		}
		rp.ServeHTTP(w, req)
	}
}

// PickUpstream is exposed for the UDP layer.
func (r *Router) PickUpstream(ctx context.Context, name string) (*url.URL, error) {
	return r.pickUpstream(ctx, name)
}

package registry

import (
	"context"
	"sync"
	"time"

	"github.com/vpt/common/httpx"
)

type ServiceInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"` // host:port
	Scheme  string `json:"scheme"`  // http
}

type RegisterResp struct {
	InstanceID string `json:"instance_id"`
}

type Client struct {
	hc         *httpx.Client
	instanceID string
	mu         sync.RWMutex
	cache      map[string][]ServiceInfo
	config     map[string]string
}

func NewClient(registryURL string) *Client {
	return &Client{
		hc:     httpx.NewClient(registryURL),
		cache:  make(map[string][]ServiceInfo),
		config: make(map[string]string),
	}
}

func (c *Client) Register(ctx context.Context, info ServiceInfo) error {
	var resp RegisterResp
	if err := c.hc.PostJSON(ctx, "/api/v1/registry/register", info, &resp); err != nil {
		return err
	}
	c.mu.Lock()
	c.instanceID = resp.InstanceID
	c.mu.Unlock()
	return nil
}

func (c *Client) Heartbeat(ctx context.Context) error {
	c.mu.RLock()
	id := c.instanceID
	c.mu.RUnlock()
	return c.hc.PostJSON(ctx, "/api/v1/registry/heartbeat", map[string]string{"instance_id": id}, nil)
}

func (c *Client) Discover(ctx context.Context, name string) ([]ServiceInfo, error) {
	var list []ServiceInfo
	if err := c.hc.GetJSON(ctx, "/api/v1/registry/discover?name="+name, &list); err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.cache[name] = list
	c.mu.Unlock()
	return list, nil
}

func (c *Client) GetConfig(ctx context.Context, key string) (string, error) {
	var resp struct {
		Value string `json:"value"`
	}
	if err := c.hc.GetJSON(ctx, "/api/v1/config/get?key="+key, &resp); err != nil {
		return "", err
	}
	c.mu.Lock()
	c.config[key] = resp.Value
	c.mu.Unlock()
	return resp.Value, nil
}

func (c *Client) StartHeartbeat(ctx context.Context, interval time.Duration) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				_ = c.Heartbeat(ctx)
			}
		}
	}()
}

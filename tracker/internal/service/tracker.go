package service

import (
	"context"
	"errors"
	"time"

	"github.com/vpt/common/httpx"
	"github.com/vpt/common/registry"
	"github.com/vpt/tracker/internal/config"
	"github.com/vpt/tracker/internal/model"
	"github.com/vpt/tracker/internal/repo"
)

type Tracker struct {
	peers *repo.PeerRepo
	stats *repo.StatRepo
	reg   *registry.Client
	cfg   *config.Config
}

func NewTracker(p *repo.PeerRepo, s *repo.StatRepo, reg *registry.Client, cfg *config.Config) *Tracker {
	return &Tracker{peers: p, stats: s, reg: reg, cfg: cfg}
}

// Announce processes an announce request: authenticate via passkey, upsert peer, return peer list.
func (t *Tracker) Announce(ctx context.Context, req *model.AnnounceRequest) (*model.AnnounceResponse, error) {
	userID, err := t.lookupPasskey(ctx, req.Passkey)
	if err != nil {
		return nil, err
	}

	if req.Event == "stopped" {
		_ = t.peers.Delete(req.InfoHash, req.PeerID)
	} else {
		p := &model.Peer{
			InfoHash: req.InfoHash, PeerID: req.PeerID, UserID: userID,
			IP: req.IP, Port: req.Port,
			Uploaded: req.Uploaded, Downloaded: req.Downloaded, Left: req.Left,
			Event: req.Event, LastSeen: time.Now(),
		}
		if err := t.peers.Upsert(p); err != nil {
			return nil, err
		}
	}

	seeders, leechers, _ := t.peers.CountSeedLeech(req.InfoHash)
	completedDelta := int64(0)
	if req.Event == "completed" {
		completedDelta = 1
	}
	_ = t.stats.Update(req.InfoHash, seeders, leechers, completedDelta)

	numWant := req.NumWant
	if numWant <= 0 || numWant > 50 {
		numWant = 50
	}
	peers, _ := t.peers.ListPeers(req.InfoHash, numWant)

	return &model.AnnounceResponse{
		Interval:   t.cfg.AnnounceInterval,
		Complete:   seeders,
		Incomplete: leechers,
		Peers:      peers,
	}, nil
}

func (t *Tracker) Scrape(infoHash string) (*model.StatSnapshot, error) {
	return t.stats.Get(infoHash)
}

func (t *Tracker) ListStats(limit, offset int) ([]model.StatSnapshot, error) {
	return t.stats.List(limit, offset)
}

func (t *Tracker) StartReaper(ctx context.Context, interval time.Duration) {
	tk := time.NewTicker(interval)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			_ = t.peers.PurgeStale(time.Now().Add(-time.Duration(t.cfg.AnnounceInterval*2) * time.Second))
		}
	}
}

// lookupPasskey calls admin /api/v1/auth/passkey to translate passkey -> user_id.
func (t *Tracker) lookupPasskey(ctx context.Context, passkey string) (int64, error) {
	if passkey == "" {
		return 0, errors.New("missing passkey")
	}
	instances, err := t.reg.Discover(ctx, "admin")
	if err != nil || len(instances) == 0 {
		return 0, errors.New("admin unavailable")
	}
	base := instances[0].Scheme + "://" + instances[0].Address
	hc := httpx.NewClient(base)
	var resp struct {
		Code int `json:"code"`
		Data struct {
			UserID int64 `json:"user_id"`
		} `json:"data"`
	}
	if err := hc.GetJSON(ctx, "/api/v1/auth/passkey?passkey="+passkey, &resp); err != nil {
		return 0, err
	}
	if resp.Data.UserID == 0 {
		return 0, errors.New("invalid passkey")
	}
	return resp.Data.UserID, nil
}

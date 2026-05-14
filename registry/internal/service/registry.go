package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/vpt/registry/internal/model"
	"github.com/vpt/registry/internal/repo"
)

type Registry struct {
	repo *repo.ServiceRepo
}

func NewRegistry(r *repo.ServiceRepo) *Registry { return &Registry{repo: r} }

func (s *Registry) Register(name, address, scheme string) (*model.ServiceInstance, error) {
	id := newID()
	inst := &model.ServiceInstance{
		InstanceID: id, Name: name, Address: address, Scheme: scheme,
		Status: "up", LastSeen: time.Now(),
	}
	if err := s.repo.Upsert(inst); err != nil {
		return nil, err
	}
	return inst, nil
}

func (s *Registry) Heartbeat(instanceID string) error {
	return s.repo.Touch(instanceID, time.Now())
}

func (s *Registry) Discover(name string) ([]model.ServiceInstance, error) {
	return s.repo.ListByName(name)
}

func (s *Registry) StartHealthCheck(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			_ = s.repo.MarkDownBefore(time.Now().Add(-30 * time.Second))
		}
	}
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

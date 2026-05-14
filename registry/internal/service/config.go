package service

import (
	"github.com/vpt/registry/internal/model"
	"github.com/vpt/registry/internal/repo"
)

type Config struct {
	repo *repo.ConfigRepo
}

func NewConfig(r *repo.ConfigRepo) *Config { return &Config{repo: r} }

func (s *Config) Get(key string) (*model.ConfigItem, error) { return s.repo.Get(key) }
func (s *Config) Set(key, value string) error              { return s.repo.Set(key, value) }
func (s *Config) List() ([]model.ConfigItem, error)        { return s.repo.List() }

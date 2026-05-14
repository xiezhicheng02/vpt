package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/vpt/admin/internal/config"
	"github.com/vpt/admin/internal/model"
	"github.com/vpt/admin/internal/repo"
)

type Torrent struct {
	repo *repo.TorrentRepo
	cfg  *config.Config
}

func NewTorrent(r *repo.TorrentRepo, cfg *config.Config) *Torrent {
	return &Torrent{repo: r, cfg: cfg}
}

// Upload saves the torrent file to disk and creates a DB record.
// TODO: parse the .torrent file (bencode) to extract info_hash, name, size.
func (s *Torrent) Upload(uploaderID int64, filename string, body io.Reader) (*model.Torrent, error) {
	if filename == "" {
		return nil, errors.New("missing filename")
	}
	if err := os.MkdirAll(s.cfg.TorrentDir, 0o755); err != nil {
		return nil, err
	}
	dst := filepath.Join(s.cfg.TorrentDir, filename)
	f, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(f, body); err != nil {
		return nil, err
	}
	t := &model.Torrent{
		InfoHash:   "", // TODO: compute from bencoded info dict
		Name:       filename,
		UploaderID: uploaderID,
		FilePath:   dst,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Torrent) List(limit, offset int) ([]model.Torrent, error) {
	return s.repo.List(limit, offset)
}

func (s *Torrent) Get(infoHash string) (*model.Torrent, error) {
	return s.repo.FindByInfoHash(infoHash)
}

func (s *Torrent) Delete(id int64) error {
	return s.repo.Delete(id)
}

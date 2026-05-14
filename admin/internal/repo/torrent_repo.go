package repo

import (
	"database/sql"

	"github.com/vpt/admin/internal/model"
)

type TorrentRepo struct{ db *sql.DB }

func NewTorrentRepo(db *sql.DB) *TorrentRepo { return &TorrentRepo{db: db} }

func (r *TorrentRepo) Create(t *model.Torrent) error {
	res, err := r.db.Exec(`INSERT INTO torrents(info_hash,name,size,uploader_id,file_path,created_at)
		VALUES(?,?,?,?,?,?)`, t.InfoHash, t.Name, t.Size, t.UploaderID, t.FilePath, t.CreatedAt)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	t.ID = id
	return nil
}

func (r *TorrentRepo) FindByInfoHash(hash string) (*model.Torrent, error) {
	row := r.db.QueryRow(`SELECT id,info_hash,name,size,uploader_id,file_path,created_at FROM torrents WHERE info_hash=?`, hash)
	var t model.Torrent
	if err := row.Scan(&t.ID, &t.InfoHash, &t.Name, &t.Size, &t.UploaderID, &t.FilePath, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TorrentRepo) List(limit, offset int) ([]model.Torrent, error) {
	rows, err := r.db.Query(`SELECT id,info_hash,name,size,uploader_id,file_path,created_at FROM torrents ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Torrent
	for rows.Next() {
		var t model.Torrent
		if err := rows.Scan(&t.ID, &t.InfoHash, &t.Name, &t.Size, &t.UploaderID, &t.FilePath, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *TorrentRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM torrents WHERE id=?`, id)
	return err
}

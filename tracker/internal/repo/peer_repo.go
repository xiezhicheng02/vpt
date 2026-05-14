package repo

import (
	"database/sql"
	"time"

	"github.com/vpt/tracker/internal/model"
)

type PeerRepo struct{ db *sql.DB }

func NewPeerRepo(db *sql.DB) *PeerRepo { return &PeerRepo{db: db} }

func (r *PeerRepo) Upsert(p *model.Peer) error {
	_, err := r.db.Exec(`INSERT INTO peers(info_hash,peer_id,user_id,ip,port,uploaded,downloaded,left_bytes,event,last_seen)
		VALUES(?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(info_hash, peer_id) DO UPDATE SET
		  user_id=excluded.user_id, ip=excluded.ip, port=excluded.port,
		  uploaded=excluded.uploaded, downloaded=excluded.downloaded,
		  left_bytes=excluded.left_bytes, event=excluded.event, last_seen=excluded.last_seen`,
		p.InfoHash, p.PeerID, p.UserID, p.IP, p.Port, p.Uploaded, p.Downloaded, p.Left, p.Event, p.LastSeen)
	return err
}

func (r *PeerRepo) Delete(infoHash, peerID string) error {
	_, err := r.db.Exec(`DELETE FROM peers WHERE info_hash=? AND peer_id=?`, infoHash, peerID)
	return err
}

func (r *PeerRepo) ListPeers(infoHash string, limit int) ([]model.PeerAddr, error) {
	rows, err := r.db.Query(`SELECT ip, port FROM peers WHERE info_hash=? LIMIT ?`, infoHash, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.PeerAddr
	for rows.Next() {
		var p model.PeerAddr
		if err := rows.Scan(&p.IP, &p.Port); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *PeerRepo) CountSeedLeech(infoHash string) (seeders, leechers int, err error) {
	row := r.db.QueryRow(`SELECT
		SUM(CASE WHEN left_bytes=0 THEN 1 ELSE 0 END),
		SUM(CASE WHEN left_bytes>0 THEN 1 ELSE 0 END)
		FROM peers WHERE info_hash=?`, infoHash)
	var s, l sql.NullInt64
	if err = row.Scan(&s, &l); err != nil {
		return 0, 0, err
	}
	return int(s.Int64), int(l.Int64), nil
}

func (r *PeerRepo) PurgeStale(before time.Time) error {
	_, err := r.db.Exec(`DELETE FROM peers WHERE last_seen < ?`, before)
	return err
}

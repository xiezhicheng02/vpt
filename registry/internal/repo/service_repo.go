package repo

import (
	"database/sql"
	"time"

	"github.com/vpt/registry/internal/model"
)

type ServiceRepo struct{ db *sql.DB }

func NewServiceRepo(db *sql.DB) *ServiceRepo { return &ServiceRepo{db: db} }

func (r *ServiceRepo) Upsert(s *model.ServiceInstance) error {
	_, err := r.db.Exec(`INSERT INTO services(instance_id,name,address,scheme,status,last_seen)
		VALUES(?,?,?,?,?,?)
		ON CONFLICT(instance_id) DO UPDATE SET name=excluded.name,address=excluded.address,
		scheme=excluded.scheme,status=excluded.status,last_seen=excluded.last_seen`,
		s.InstanceID, s.Name, s.Address, s.Scheme, s.Status, s.LastSeen)
	return err
}

func (r *ServiceRepo) Touch(instanceID string, t time.Time) error {
	_, err := r.db.Exec(`UPDATE services SET last_seen=?, status='up' WHERE instance_id=?`, t, instanceID)
	return err
}

func (r *ServiceRepo) ListByName(name string) ([]model.ServiceInstance, error) {
	rows, err := r.db.Query(`SELECT instance_id,name,address,scheme,status,last_seen FROM services WHERE name=? AND status='up'`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.ServiceInstance
	for rows.Next() {
		var s model.ServiceInstance
		if err := rows.Scan(&s.InstanceID, &s.Name, &s.Address, &s.Scheme, &s.Status, &s.LastSeen); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *ServiceRepo) MarkDownBefore(t time.Time) error {
	_, err := r.db.Exec(`UPDATE services SET status='down' WHERE last_seen < ?`, t)
	return err
}

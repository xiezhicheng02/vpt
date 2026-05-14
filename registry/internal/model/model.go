package model

import "time"

type ServiceInstance struct {
	InstanceID string    `json:"instance_id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Scheme     string    `json:"scheme"`
	Status     string    `json:"status"` // up / down
	LastSeen   time.Time `json:"last_seen"`
}

type ConfigItem struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

package model

import "time"

type Peer struct {
	ID         int64
	InfoHash   string
	PeerID     string
	UserID     int64
	IP         string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string // started / completed / stopped / ""
	LastSeen   time.Time
}

type AnnounceRequest struct {
	InfoHash   string
	PeerID     string
	Passkey    string
	IP         string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
	NumWant    int
	Compact    bool
}

type AnnounceResponse struct {
	Interval   int
	Complete   int // seeders
	Incomplete int // leechers
	Peers      []PeerAddr
}

type PeerAddr struct {
	IP   string
	Port int
}

type StatSnapshot struct {
	InfoHash   string `json:"info_hash"`
	Seeders    int    `json:"seeders"`
	Leechers   int    `json:"leechers"`
	Completed  int64  `json:"completed"`
	UpdatedAt  time.Time `json:"updated_at"`
}

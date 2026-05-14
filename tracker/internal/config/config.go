package config

import "os"

type Config struct {
	HTTPAddr     string
	UDPAddr      string
	DBPath       string
	RegistryURL  string
	AnnounceInterval int // seconds
}

func Load() *Config {
	return &Config{
		HTTPAddr:         getenv("TRACKER_HTTP_ADDR", ":8002"),
		UDPAddr:          getenv("TRACKER_UDP_ADDR", ":8003"),
		DBPath:           getenv("TRACKER_DB", "data/tracker.db"),
		RegistryURL:      getenv("REGISTRY_URL", "http://127.0.0.1:8500"),
		AnnounceInterval: 1800,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

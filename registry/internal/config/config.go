package config

import "os"

type Config struct {
	Addr   string
	DBPath string
}

func Load() *Config {
	return &Config{
		Addr:   getenv("REGISTRY_ADDR", ":8500"),
		DBPath: getenv("REGISTRY_DB", "data/registry.db"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

package config

import "os"

type Config struct {
	HTTPAddr    string
	UDPAddr     string
	RegistryURL string
}

func Load() *Config {
	return &Config{
		HTTPAddr:    getenv("GATEWAY_HTTP_ADDR", ":8000"),
		UDPAddr:     getenv("GATEWAY_UDP_ADDR", ":8001"),
		RegistryURL: getenv("REGISTRY_URL", "http://127.0.0.1:8500"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

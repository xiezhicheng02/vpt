package config

import "os"

type Config struct {
	Addr         string
	DBPath       string
	RegistryURL  string
	TorrentDir   string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	MailFrom     string
	TokenSecret  string
}

func Load() *Config {
	return &Config{
		Addr:        getenv("ADMIN_ADDR", ":8001"),
		DBPath:      getenv("ADMIN_DB", "data/admin.db"),
		RegistryURL: getenv("REGISTRY_URL", "http://127.0.0.1:8500"),
		TorrentDir:  getenv("TORRENT_DIR", "data/torrents"),
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    getenv("SMTP_PORT", "587"),
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASS"),
		MailFrom:    getenv("MAIL_FROM", "no-reply@vpt.local"),
		TokenSecret: getenv("TOKEN_SECRET", "please-change-me"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

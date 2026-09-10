package config

import ()

type Config struct {
	APIBase      string
	Player       string
	Quality      string
	Radio        bool
	QueueSize    int
	HistoryLimit int
	CookieMode   string
	CookiesFile  string
	Debug        bool
	Volume       int
}

func Load() *Config {
	// For simplicity in this demo, just return defaults.
	// Production should parse ~/.config/soundcloud-radio/config.toml
	return &Config{
		APIBase:      "https://api-v2.soundcloud.com",
		Player:       "mpv",
		Quality:      "best",
		Radio:        true,
		QueueSize:    10,
		HistoryLimit: 5000,
		CookieMode:   "auto",
		Debug:        false,
		Volume:       100,
	}
}

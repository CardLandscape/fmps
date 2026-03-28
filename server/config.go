package main

import (
	"os"
	"path/filepath"
)

type Config struct {
	DBPath    string
	Port      string
	JWTSecret string
}

func LoadConfig() Config {
	cfg := Config{
		Port:      "8080",
		JWTSecret: "fmps-secret-key-2024",
	}

	if p := os.Getenv("FMPS_PORT"); p != "" {
		cfg.Port = p
	}
	if s := os.Getenv("FMPS_JWT_SECRET"); s != "" {
		cfg.JWTSecret = s
	}

	if dbPath := os.Getenv("FMPS_DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		dbDir := filepath.Join(homeDir, ".fmps")
		_ = os.MkdirAll(dbDir, 0755)
		cfg.DBPath = filepath.Join(dbDir, "fmps.db")
	}

	return cfg
}

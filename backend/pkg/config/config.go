package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port          string
	DatabaseURL   string
	GCPProjectID  string
	GCPRegion     string
}

// Load reads required environment variables and returns a Config.
// It fails fast if any required variable is missing.
func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnvOrDefault("PORT", "8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		GCPProjectID: os.Getenv("GCP_PROJECT_ID"),
		GCPRegion:    getEnvOrDefault("GCP_REGION", "asia-northeast1"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("required env var DATABASE_URL is not set")
	}
	if cfg.GCPProjectID == "" {
		return nil, fmt.Errorf("required env var GCP_PROJECT_ID is not set")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

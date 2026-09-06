// Package config loads and validates the wishlist API's environment-based
// configuration. Importing it triggers autoloading of a .env file (if
// present) via github.com/joho/godotenv/autoload.
package config

import (
	"errors"
	"os"
	"strings"
	
	// this will automatically load your .env file:
	_ "github.com/joho/godotenv/autoload"
)

// Config holds all runtime configuration for the wishlist API server.
type Config struct {
	Port            string
	DatabaseURL     string
	SupabaseURL     string
	SupabaseJWKSURL string
}

// LoadConfig reads configuration from environment variables (populated from
// .env by the godotenv autoload import above, if present). It returns an
// error if any required variable is missing.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		SupabaseURL:     os.Getenv("SUPABASE_URL"),
		SupabaseJWKSURL: os.Getenv("SUPABASE_JWKS_URL"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	if cfg.SupabaseJWKSURL == "" {
		if cfg.SupabaseURL == "" {
			return nil, errors.New("SUPABASE_JWKS_URL or SUPABASE_URL environment variable is required")
		}
		cfg.SupabaseJWKSURL = strings.TrimRight(cfg.SupabaseURL, "/") + "/auth/v1/.well-known/jwks.json"
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

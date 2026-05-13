package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config holds application configuration.
type Config struct {
	Port          string
	PostgresDSN   string
	ClickHouseDSN string
	KafkaBrokers  []string
	JWTSecret     string
}

// Load reads configuration from env.
func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Port:        getEnv("PORT", "8080"),
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/historian?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-change-me"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

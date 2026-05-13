package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
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
		Port:          getEnv("PORT", "8080"),
		PostgresDSN:   getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/historian?sslmode=disable"),
		ClickHouseDSN: getEnv("CLICKHOUSE_DSN", "clickhouse://default:clickhouse_secret@localhost:9000/analytics"),
		KafkaBrokers:  strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-change-me"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

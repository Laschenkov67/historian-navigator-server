package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates connection pool.
func NewPostgresPool(dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 20
	return pgxpool.NewWithConfig(context.Background(), cfg)
}

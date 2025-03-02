package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenPool(ctx context.Context, config DatabaseConfig) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(config.ConnStr)
	if err != nil {
		return nil, err
	}
	pgxCfg.MaxConnIdleTime = config.IdleTime()
	pgxCfg.MaxConns = config.MaxOpenConns

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

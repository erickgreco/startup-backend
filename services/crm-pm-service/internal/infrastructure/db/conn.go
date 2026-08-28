package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	MaxConns              int32
	MinConns              int32
	MaxConnLifeTime       time.Duration
	MaxConnIdleTime       time.Duration
	HealthCheckPeriod     time.Duration
	MaxConnLifeTimeJitter time.Duration
}

func Connect(conn string, cfg Config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = cfg.MaxConns
	conf.MinConns = cfg.MinConns
	conf.MaxConnLifetime = cfg.MaxConnLifeTime
	conf.MaxConnIdleTime = cfg.MaxConnIdleTime
	conf.HealthCheckPeriod = cfg.HealthCheckPeriod
	conf.MaxConnLifetimeJitter = cfg.MaxConnLifeTimeJitter

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

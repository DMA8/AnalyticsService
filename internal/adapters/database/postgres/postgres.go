// Package postgres implements postgres connection.
package postgres

import (
	"gitlab.com/g6834/team31/analytics/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database -.
type Database struct {
	Pool *pgxpool.Pool
}

// New -.
func New(ctx context.Context, cfg config.PG) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &Database{
		Pool: pool,
	}, nil
}

// Close -.
func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

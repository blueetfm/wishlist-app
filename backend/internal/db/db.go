// Package db handles database connectivity and query helpers for the
// wishlist API.
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a pgx connection pool for the given Postgres connection
// string (e.g. a Supabase DATABASE_URL) and verifies connectivity with a
// ping before returning.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

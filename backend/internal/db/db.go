// Package db handles database connectivity and query helpers for the
// wishlist API.
package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Use a Pool interface
// *pgxpool.Pool can easily implement this interface, but importantly,
// this allows unit testing to substitute with a mock 
// without needing a real, live Postgres connection.
type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

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

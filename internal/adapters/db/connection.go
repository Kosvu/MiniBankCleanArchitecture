package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, "postgres://postgres:usif775shakh@localhost:5432/minibank?sslmode=disable")

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, err
}

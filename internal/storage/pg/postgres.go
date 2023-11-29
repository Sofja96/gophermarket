package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewStorage(ctx context.Context, dsn string) (*Postgres, error) {
	err := migrateDatabase(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable db migrate: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create pgxpool: %w", err)
	}
	dbc := &Postgres{
		DB: pool,
	}
	return dbc, nil
}

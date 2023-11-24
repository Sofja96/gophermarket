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
	//err := migrateDatabase(dsn)
	//if err != nil {
	//	return nil, fmt.Errorf("unable db migrate: %w", err)
	//}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create pgxpool: %w", err)
	}
	dbc := &Postgres{
		DB: pool,
	}
	err = dbc.init(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error init db: %w", err)
	}
	return dbc, nil
}

func (pg *Postgres) init(ctx context.Context) error {
	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	tx.Exec(ctx, `
		CREATE TABLE "users" (
   "id" bigserial PRIMARY KEY,
   "login" varchar NOT NULL,
   "balance" real NOT NULL default 0,
   "withdrawn" real NOT NULL default 0,
   "password" varchar NOT NULL
);`)

	tx.Exec(ctx, `
CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES users(id),
  "number" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT ('NEW'),
  "accrual" real NOT NULL default 0,
  "uploaded_at" timestamptz NOT NULL DEFAULT (now())
);
	`)

	tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS balance(
			id SERIAL PRIMARY KEY,
			current NUMERIC(10,2),
			withdrawn NUMERIC(10,2),
			fk_user_id INTEGER REFERENCES users(id) NOT NULL
		)
	`)

	tx.Exec(ctx, `
CREATE TABLE "withdrawals" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES users(id),
  "number" varchar NOT NULL,
  "sum" real NOT NULL,
  "processed_at"  timestamptz NOT NULL DEFAULT (now())
);
	`)

	return tx.Commit(ctx)
}

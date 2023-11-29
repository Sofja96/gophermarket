package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) CreateUser(user string, password string) (string, error) {
	ctx := context.Background()
	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, "INSERT INTO users (login, password) VALUES ($1, $2)", user, password)
	if err != nil {
		return "", fmt.Errorf("error insert user: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}
	return user, nil
}

func (pg *Postgres) GetUserHashPassword(user string) (string, error) {
	ctx := context.Background()
	var pass string
	row := pg.DB.QueryRow(ctx, "SELECT password FROM users WHERE login = $1", user)
	err := row.Scan(&pass)
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", fmt.Errorf("unable select pass: %w", err)
		} else {
			return "", nil
		}
	}

	return pass, nil
}

func (pg *Postgres) GetUserIDByName(user string) (bool, error) {
	ctx := context.Background()
	var id string
	row := pg.DB.QueryRow(ctx, "SELECT id FROM users WHERE login = $1", user)
	err := row.Scan(&id)
	if err != nil {
		if err != pgx.ErrNoRows {
			return false, fmt.Errorf("unable select id: %w", err)
		} else {
			return false, nil
		}
	}

	return true, nil
}

func (pg *Postgres) GetUserID(user string) (string, error) {
	ctx := context.Background()
	var id string
	row := pg.DB.QueryRow(ctx, "SELECT id FROM users WHERE login = $1", user)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("unable select id: %w", err)
	}

	return id, nil
}

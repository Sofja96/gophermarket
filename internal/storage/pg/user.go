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

//func (pg *Postgres) CreateUser(user, password string) (string, error) {
//	ctx := context.Background()
//	tx, err := pg.DB.Begin(ctx)
//	if err != nil {
//		return "", err
//	}
//	defer func() { _ = tx.Rollback(ctx) }()
//	//var id string
//	//raw := tx.Query(ctx, "SELECT value FROM counter_metrics WHERE name = $1", id)
//	_, err = tx.Query(ctx, `INSERT INTO users (login, password) VALUES ($1, $2)`, user, password)
//	//raw, err := tx.Query(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", user, password)
//	if err != nil {
//		return "", err
//	}
//	//err = raw.Scan(&id)
//	//if err != nil {
//	//	return err
//	//}
//	err = tx.Commit(ctx)
//	if err != nil {
//		return "", err
//	}
//	return user, nil
//
//}

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

//func (pg *Postgres) getUserID(user string) (string, error) {
//	ctx := context.Background()
//	var userID string
//	row := pg.DB.QueryRow(ctx, `SELECT id FROM users WHERE name = $1`, user)
//	if err := row.Scan(&userID); err != nil {
//		return userID, fmt.Errorf("error on scanning values: %w", err)
//	}
//	return userID, nil
//}

//func (pg *Postgres) CreateUser(ctx context.Context, user, password, hash, salt string) error {
//	tx, err := pg.DB.Begin(ctx)
//	if err != nil {
//		return err
//	}
//	defer func() { _ = tx.Rollback(ctx) }()
//	//	var id uint
//	_, err = tx.Exec(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", user, password)
//	//	err = tx.Scan(&id)
//	if err != nil {
//		return err
//	}
//
//	err = tx.Commit(ctx)
//	if err != nil {
//		return err
//	}
//
//	tx.Commit(ctx)
//	return nil
//
//}

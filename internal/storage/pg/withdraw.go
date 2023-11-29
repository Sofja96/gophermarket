package pg

import (
	"context"
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/gommon/log"
)

func (pg *Postgres) WithdrawBalance(user, orderNumber string, sum float32) error {
	ctx := context.Background()
	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	userID, err := pg.GetUserID(user)
	if err != nil {
		return fmt.Errorf("error get id from users: %w", err)
	}
	log.Print(userID)

	balance, err := pg.GetBalance(user)
	if err != nil {
		return fmt.Errorf("error get balance from users: %w", err)
	}

	if balance.Current < sum {
		return helpers.ErrInsufficientBalance
	}

	newBalance := balance.Current - sum

	_, err = tx.Exec(ctx, "UPDATE users SET balance = $1, withdrawn = withdrawn + $2 WHERE id = $3", newBalance, sum, userID)
	if err != nil {
		log.Infof("error update values in users with balance")
		return fmt.Errorf("error update values in users with balance: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO withdrawals (user_id, number, sum) VALUES ($1, $2, $3)", userID, orderNumber, sum)
	if err != nil {
		log.Infof("error insert values in withdrawals")
		return fmt.Errorf("error insert values in withdrawals: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	log.Infof("WithdrawBalance: Successfully updated")

	return nil
}

func (pg *Postgres) Getwithdrawals(user string) ([]models.UserWithdrawal, error) {
	ctx := context.Background()
	withdrawals := make([]models.UserWithdrawal, 0)

	userID, err := pg.GetUserID(user)
	if err != nil {
		return withdrawals, fmt.Errorf("error get if from users: %w", err)
	}

	row, err := pg.DB.Query(ctx, "SELECT number,sum,processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC", userID)
	if err != nil {
		return withdrawals, fmt.Errorf("error select values in withdrawals: %w", err)
	}
	defer row.Close()

	for row.Next() {
		var wd models.UserWithdrawal
		err := row.Scan(&wd.Order, &wd.Sum, &wd.ProcessedAt)
		if err != nil {
			return withdrawals, fmt.Errorf("error scan values in withdrawals: %w", err)
		}
		withdrawals = append(withdrawals, wd)
	}

	if err := row.Err(); err != nil {
		if err == pgx.ErrNoRows {
			return withdrawals, nil
		}

		return withdrawals, fmt.Errorf("error select values in withdrawals: %w", err)
	}

	return withdrawals, nil
}

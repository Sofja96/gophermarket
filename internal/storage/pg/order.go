package pg

import (
	"context"
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/gommon/log"
	"strings"
	"time"
)

func (pg *Postgres) CreateOrder(orderNumber, user string) (*models.Order, error) {
	ctx := context.Background()
	cctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	tx, err := pg.DB.Begin(cctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(cctx) }()

	userID, err := pg.GetUserID(user)
	if err != nil {
		return nil, fmt.Errorf("error get if from users: %w", err)
	}
	log.Print(userID)
	var orderUserID string
	row := tx.QueryRow(cctx, "SELECT user_id FROM orders WHERE number = $1", orderNumber)
	if err := row.Scan(&orderUserID); err == nil {
		if orderUserID == userID {
			log.Infof("order number already exists for this user")
			return nil, helpers.ErrExistsOrder
		}
		log.Infof("order number already exists for another user")
		return nil, helpers.ErrAnotherUserOrder
	}
	log.Print(orderUserID)

	_, err = tx.Exec(cctx, "INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)", orderNumber, userID, models.NEW)
	if err != nil {
		log.Infof("error insert")
		return nil, fmt.Errorf("error get user id: %w", err)
	}

	err = tx.Commit(cctx)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (pg *Postgres) UpdateOrder(orderNumber, status string, accrual float32) error {
	ctx := context.Background()
	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var userID uint
	row := tx.QueryRow(ctx, "SELECT user_id FROM orders WHERE number = $1", orderNumber)
	if err := row.Scan(&userID); err != nil {
		if err == pgx.ErrNoRows {
			return err // Order not found
		}
		return err // Other error occurred
	}
	log.Print(userID)

	var order models.Order
	rows := tx.QueryRow(ctx, "SELECT number,status,accrual,uploaded_at,user_id FROM orders WHERE number = $1 FOR UPDATE SKIP LOCKED", orderNumber)
	if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &userID); err != nil {
		if err == pgx.ErrNoRows {
			helpers.Infof("error scan values in orders: %s", err)
			return fmt.Errorf("error scan values in orders: %w", err)
		}
		return err
	}

	if order.Status == models.PROCESSED || order.Status == models.INVALID {
		return pgx.ErrNoRows
	}

	helpers.Infof(orderNumber, "orderNumber before update")
	_, err = tx.Exec(ctx, "UPDATE orders SET status = $1, accrual = $2 WHERE number = $3", status, accrual, orderNumber)
	if err != nil {
		//helpers.Infof("error update values in orders %s", status, accrual, orderNumber)
		//log.Infof("error update values in orders")
		return fmt.Errorf("error update orders: %w", err)
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = coalesce(balance, 0) + $1 WHERE id = $2", accrual, userID)
	if err != nil {
		//helpers.Infof("error update values in users %f\n", accrual, userID)
		return fmt.Errorf("error update users: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetOrders(user string) ([]models.Order, error) {
	ctx := context.Background()
	orders := make([]models.Order, 0)

	userID, err := pg.GetUserID(user)
	if err != nil {
		return orders, fmt.Errorf("error get if from users: %w", err)
	}

	row, err := pg.DB.Query(ctx, "SELECT number,status,accrual,uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return orders, fmt.Errorf("error select values in orders: %w", err)
	}
	defer row.Close()

	for row.Next() {
		var order models.Order
		err := row.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return orders, fmt.Errorf("error scan values in orders: %w", err)
		}
		orders = append(orders, order)
	}

	if err := row.Err(); err != nil {
		if err == pgx.ErrNoRows {
			return orders, nil
		}

		return orders, fmt.Errorf("error select values in orders: %w", err)
	}

	return orders, nil
}

func (pg *Postgres) GetBalance(user string) (models.UserBalance, error) {
	ctx := context.Background()
	var balance models.UserBalance

	row := pg.DB.QueryRow(ctx, "SELECT balance, withdrawn FROM users WHERE login = $1", user)
	if err := row.Scan(&balance.Current, &balance.Withdrawn); err != nil {
		if err == pgx.ErrNoRows {
			return balance, fmt.Errorf("error select balance in users: %w", err)
		}
		return balance, err // Other error occurred
	}
	log.Print(balance, "balance in bd")

	return balance, nil
}

func (pg *Postgres) lock(ctx context.Context, ID string, what string) error {
	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Prepare(ctx, "my-query", "SELECT id FROM "+what+" WHERE id = $1 FOR UPDATE")
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, "my-query", ID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

func (pg *Postgres) GetOrderStatus(status []string) ([]string, error) {
	ctx := context.Background()
	orders := make([]string, 0)
	statuses := strings.Join(status, ",")

	row, err := pg.DB.Query(ctx, "SELECT number FROM orders where status IN ($1)", statuses)
	if err != nil {
		return orders, fmt.Errorf("error select status in orders: %w", err)
	}
	defer row.Close()

	for row.Next() {
		var order string
		err := row.Scan(&order)
		if err != nil {
			return orders, fmt.Errorf("error scan values: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

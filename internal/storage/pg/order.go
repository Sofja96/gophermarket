package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/labstack/gommon/log"
	"strconv"
	"time"
)

func (pg *Postgres) CreateOrder(orderNumber, user string) (*models.Order, error) {
	ctx := context.Background()
	cctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	tx, err := pg.DB.Begin(cctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(cctx) }()

	var userID uint
	row := tx.QueryRow(cctx, "SELECT id FROM users WHERE login = $1", user)
	if err := row.Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Order not found
		}
		return nil, err // Other error occurred
	}
	log.Print(userID)

	var orderUserId uint
	row = tx.QueryRow(cctx, "SELECT user_id FROM orders WHERE number = $1", orderNumber)
	if err := row.Scan(&orderUserId); err == nil {
		if orderUserId == userID {
			log.Infof("order number already exists for this user")
			return nil, fmt.Errorf("order number already exists for this user: %w", err)
		}
		log.Infof("order number already exists for another user")
		return nil, fmt.Errorf("order number already exists for another user: %w", err)
	}
	//}
	//	}
	log.Print(orderUserId)

	//err == nil {
	//	//log.Infof(orderUserId, "orderUserId")
	//	log.Print(orderUserId, "orderUserId")
	//	if orderUserId == userID {
	//		return nil, fmt.Errorf("user already ordered this order: %w", err)
	//	}
	//	return nil, fmt.Errorf("order number already exists: %w", err)
	//}

	//if userID == orderUserId {
	//	return nil, fmt.Errorf("user already ordered this order: %w", err)
	//}
	//
	//if len(orderNumber) != 0 {
	//	return nil, fmt.Errorf("order number already exists: %w", err)
	//}
	//_ = tx.QueryRow(ctx, "SELECT number FROM orders WHERE number = $1", orderNumber)
	//if err != nil {
	//	return nil, fmt.Errorf("order number already exists: %w", err)
	//}
	//ok := tx.QueryRow(ctx, "SELECT number = $1 FROM orders", orderNumber)
	//if !ok {
	//	return nil, fmt.Errorf("order number already exists: %w", err)
	//}

	//err != nil {
	//	return nil, fmt.Errorf("order number already exists: %w", err)
	//}

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

	//defer smt.Close()

	//_, err = smt.ExecContext(ctx, ID)
	_, err = tx.Exec(ctx, "my-query", ID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

func (pg *Postgres) LockUsers(ctx context.Context, userID string) error {
	return pg.lock(ctx, userID, "users")
}

func (pg *Postgres) LockOrders(ctx context.Context, orderID string) error {
	return pg.lock(ctx, orderID, "orders")
}

func (pg *Postgres) UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) error {
	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = pg.LockOrders(ctx, strconv.Itoa(int(orderID)))
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Prepare(ctx, "status", "UPDATE orders SET status = $1 WHERE id = $2")
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(ctx, "status", status, orderID)
	//_, err = oStmt.ExecContext(ctx, status, orderID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return nil
}

func (pg *Postgres) UpdateOrderAccrualAndUserBalance(ctx context.Context, orderID uint, userID string, accrualResp models.OrderAccrual) error {
	log.Infof("UpdateOrderAccrualAndUserBalance params: orderID: %d, userID: %d", orderID, userID)

	tx, err := pg.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = pg.LockOrders(ctx, strconv.Itoa(int(orderID)))
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Prepare(ctx, "accrual", "UPDATE orders SET accrual = $1, status = $2 WHERE id = $3")
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	//defer oStmt.Close()

	_, err = tx.Exec(ctx, "accrual", accrualResp.Accrual, accrualResp.Status, orderID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = pg.LockUsers(ctx, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Prepare(ctx, "balance", "UPDATE users SET balance = balance + $1 WHERE id = $2")
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	//defer uStmt.Close()

	_, err = tx.Exec(ctx, "balance", accrualResp.Accrual, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	log.Infof("UpdateOrderAccrualAndUserBalance: Successfully updated")

	return nil
}

//func (pg *Postgres) UpdateOrder(orderNumber, status string, accrual float32) error {
//	ctx := context.Background()
//	tx, err := pg.DB.Begin(ctx)
//	if err != nil {
//		return err
//	}
//	defer func() { _ = tx.Rollback(ctx) }()
//
//	user, err := s.GetUserByOrder(ctx, number)
//	if err != nil {
//		return fmt.Errorf("error on getting user by order: %w", err)
//	}
//
//	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE name = $2`, accrual, user)
//	if err != nil {
//		return fmt.Errorf("error on updating user balance %w", err)
//	}
//
//	_, err = tx.Exec(
//		ctx,
//		`UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`,
//		status,
//		accrual,
//		orderNumber,
//	)
//
//	if err != nil {
//		return fmt.Errorf("error on updating values: %w", err)
//	}
//
//	if err := tx.Commit(ctx); err != nil {
//		return fmt.Errorf("error on tx commit: %w", err)
//	}
//
//	return nil
//}

//func (pg *Postgres) GetOrderIDByNumber(user string) (bool, error) {
//	ctx := context.Background()
//	var id string
//	row := pg.DB.QueryRow(ctx, "SELECT id FROM users WHERE login = $1", user)
//	err := row.Scan(&id)
//	if err != nil {
//		if err != pgx.ErrNoRows {
//			return false, fmt.Errorf("unable select id: %w", err)
//		} else {
//			return false, nil
//		}
//	}
//
//	return true, nil
//}

//func (pg *Postgres) GetOrderByNumber(orderNumber string) (*models.Order, error) {
//	ctx := context.Background()
//	tx, err := pg.DB.Begin(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer func() { _ = tx.Rollback(ctx) }()
//	var orderUserId string
//	row := tx.QueryRow(ctx, "SELECT id, user_id,number, status, accrual,uploaded_at FROM orders WHERE number = $1", orderNumber)
//	//	var order models.Order
//	if err := row.Scan(
//		//&order.Number,
//		//&order.Status,
//		//&order.Accrual,
//		//&order.UploadedAt,
//		&orderUserId,
//	); err != nil {
//		if err == sql.ErrNoRows {
//			return nil, nil // Order not found
//		}
//		return nil, err // Other error occurred
//	}
//
//	return nil, err
//}

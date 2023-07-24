package orders

import (
	"context"
	stdErr "errors"
	"github.com/jackc/pgx/v4"
	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/logger"
	"gomarket/internal/storage"
	"gomarket/internal/storage/db"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, orderID string, accountID string) *errors.ErrorApp
	GetOrderByID(ctx context.Context, orderID string) (*entities.Order, *errors.ErrorApp)
	GetOrdersByAccountID(ctx context.Context, accountID string) ([]entities.Order, *errors.ErrorApp)
	GetOrderIDsForAccrual(ctx context.Context) ([]string, *errors.ErrorApp)
	UpdateAfterAccrual(ctx context.Context, orderID string, status string, points float64) *errors.ErrorApp
}

var _ OrderRepository = &OrderPG{}

type OrderPG struct {
	pg *db.Postgres
}

func NewOrderPG(db *db.Postgres) *OrderPG {
	return &OrderPG{
		pg: db,
	}
}

func (o *OrderPG) CreateOrder(ctx context.Context, orderID string, accountID string) *errors.ErrorApp {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `INSERT INTO orders (id, account_id, status, uploaded_at)  VALUES ($1, $2, $3, $4)`
	_, err := o.pg.DB.Exec(ctx, q,
		orderID,
		accountID,
		entities.New,
		time.Now(),
	)
	if err != nil {
		if storage.IsUniqueViolation(err) {
			return errors.NewErrUniquenessViolation(err.Error())
		}

		return errors.NewErrInternal(err.Error())
	}

	return nil
}

func (o *OrderPG) GetOrderByID(ctx context.Context, orderID string) (*entities.Order, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT id, account_id, status, points  FROM orders WHERE id = $1`
	var order entities.Order
	err := o.pg.DB.QueryRow(ctx, q,
		orderID,
	).Scan(
		&order.ID,
		&order.AccountID,
		&order.Status,
		&order.Points,
	)
	if err != nil {
		if storage.IsNotFound(err) {
			return nil, errors.NewErrNotFound()
		}

		return nil, errors.NewErrInternal(err.Error())
	}

	return &order, nil
}

func (o *OrderPG) GetOrdersByAccountID(ctx context.Context, accountID string) ([]entities.Order, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT id, account_id, status, points, uploaded_at  FROM orders WHERE account_id = $1 ORDER BY uploaded_at DESC`

	rows, err := o.pg.DB.Query(ctx, q, accountID)
	if err != nil {
		if stdErr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewErrNotFound()
		}

		return nil, errors.NewErrInternal(err.Error())
	}
	defer rows.Close()

	var orders []entities.Order

	for rows.Next() {
		var order entities.Order
		err = rows.Scan(
			&order.ID,
			&order.AccountID,
			&order.Status,
			&order.Points,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, errors.NewErrInternal(err.Error())
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewErrInternal(err.Error())
	}

	return orders, nil
}

func (o *OrderPG) GetOrderIDsForAccrual(ctx context.Context) ([]string, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT id  FROM orders WHERE status IN ($1)`
	rows, err := o.pg.DB.Query(ctx, q,
		entities.New,
		//entities.Processed,
	)
	if err != nil {
		if stdErr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewErrNotFound()
		}

		return nil, errors.NewErrInternal(err.Error())
	}
	defer rows.Close()

	var orderIDs []string

	for rows.Next() {
		var id string
		err = rows.Scan(
			&id,
		)
		if err != nil {
			return nil, errors.NewErrInternal(err.Error())
		}

		orderIDs = append(orderIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewErrInternal(err.Error())
	}

	return orderIDs, nil
}

func (o *OrderPG) UpdateAfterAccrual(ctx context.Context, orderID string, status string, points float64) *errors.ErrorApp {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	tx, err := o.pg.DB.Begin(ctx)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}
	defer tx.Rollback(ctx)

	var accountID string
	q := `UPDATE orders SET status = $2, points = $3 WHERE id = $1 RETURNING account_id`
	err = tx.QueryRow(ctx, q,
		orderID,
		status,
		points,
	).Scan(&accountID)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	var balance float64
	q = `SELECT points FROM accounts WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, q,
		accountID,
	).Scan(&balance)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	logger.Log.Debug().
		Interface("balance", balance).
		Interface("points", points).
		Interface("balance+points", balance+points).
		Send()

	q = `UPDATE accounts SET points = $2 where id = $1`
	_, err = tx.Exec(
		ctx,
		q,
		accountID,
		balance+points,
	)

	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	return nil
}

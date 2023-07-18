package orders

import (
	"context"
	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/storage"
	"gomarket/internal/storage/db"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, orderID string, accountID string) *errors.ErrorApp
	GetOrdersByID(ctx context.Context, orderID string) (*entities.Order, *errors.ErrorApp)
	UpdateOrderStatus(ctx context.Context, orderID string, status int) *errors.ErrorApp
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

func (a *OrderPG) CreateOrder(ctx context.Context, orderID string, accountID string) *errors.ErrorApp {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `INSERT INTO orders (id, account_id, status, updated_at)  VALUES ($1, $2, $3, $4)`
	_, err := a.pg.DB.Exec(ctx, q,
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

func (a *OrderPG) GetOrdersByID(ctx context.Context, orderID string) (*entities.Order, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT id, account_id, status, points  FROM orders WHERE id = $1;`
	var order entities.Order
	err := a.pg.DB.QueryRow(ctx, q,
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

func (a *OrderPG) UpdateOrderStatus(ctx context.Context, orderID string, status int) *errors.ErrorApp {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `UPDATE orders SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := a.pg.DB.Exec(ctx, q,
		orderID,
		status,
		time.Now(),
	)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	return nil
}

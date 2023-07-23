package withdraw

import (
	"context"
	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/storage"
	"gomarket/internal/storage/db"
	"time"
)

type WithDrawRepository interface {
	CreateWithdraw(ctx context.Context, accountID string, orderID string, points int) *errors.ErrorApp
	GetWithdraw(ctx context.Context, accountID string) ([]entities.Withdraw, *errors.ErrorApp)
	GetWithdrawSum(ctx context.Context, accountID string) (*int, *errors.ErrorApp)
}

var _ WithDrawRepository = &WithDrawPG{}

type WithDrawPG struct {
	pg *db.Postgres
}

func NewAccountPG(db *db.Postgres) *WithDrawPG {
	return &WithDrawPG{
		pg: db,
	}
}

func (w *WithDrawPG) CreateWithdraw(ctx context.Context, accountID string, orderID string, points int) *errors.ErrorApp {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `INSERT INTO withdraws (account_id, order_id, points, processed_at)  VALUES ($1, $2, $3, $4)`
	_, err := w.pg.DB.Exec(ctx, q,
		accountID,
		orderID,
		points,
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

func (w *WithDrawPG) GetWithdraw(ctx context.Context, accountID string) ([]entities.Withdraw, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT order_id, points, processed_at  FROM withdraws WHERE account_id=$1 ORDER BY processed_at DESC`
	rows, err := w.pg.DB.Query(ctx, q,
		accountID,
	)
	if err != nil {
		return nil, errors.NewErrInternal(err.Error())
	}

	//var orID int
	defer rows.Close()
	var withdraw []entities.Withdraw
	for rows.Next() {
		var wd entities.Withdraw
		err = rows.Scan(
			&wd.OrderID,
			&wd.Sum,
			&wd.ProcessedAt,
		)
		//wd.OrderID = string(orID)
		if err == nil {
			withdraw = append(withdraw, wd)
		}
	}
	if err != nil {
		return nil, errors.NewErrInternal(err.Error())
	}

	return withdraw, nil
}

func (w *WithDrawPG) GetWithdrawSum(ctx context.Context, accountID string) (*int, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	var sum *int
	q := `SELECT SUM(points) FROM withdraws WHERE account_id=$1`
	err := w.pg.DB.QueryRow(ctx, q,
		accountID,
	).Scan(&sum)
	if err != nil {
		return nil, errors.NewErrInternal(err.Error())
	}

	return sum, nil
}

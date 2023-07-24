package accunts

import (
	"context"
	stdErr "errors"

	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/logger"
	"gomarket/internal/storage"
	"gomarket/internal/storage/db"

	"github.com/jackc/pgx/v4"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, login string, hashPass string) (string, *errors.ErrorApp)
	GetAccountByLogin(ctx context.Context, login string) (entities.Account, *errors.ErrorApp)
	GetAccountBalance(ctx context.Context, accountID string) (float64, *errors.ErrorApp)
	UpdateAccountBalance(ctx context.Context, accountID string, sum float64) *errors.ErrorApp
}

var _ AccountRepository = &AccountPG{}

type AccountPG struct {
	pg *db.Postgres
}

func NewAccountPG(db *db.Postgres) *AccountPG {
	return &AccountPG{
		pg: db,
	}
}

func (a *AccountPG) CreateAccount(ctx context.Context, login string, hashPass string) (string, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `INSERT INTO accounts (login, hash_pass)  VALUES ($1, $2) RETURNING id`
	var id string
	err := a.pg.DB.QueryRow(ctx, q,
		login,
		hashPass,
	).Scan(&id)
	if err != nil {
		if storage.IsUniqueViolation(err) {
			return "", errors.NewErrUniquenessViolation(err.Error())
		}

		return "", errors.NewErrInternal(err.Error())
	}

	return id, nil
}

func (a *AccountPG) GetAccountByLogin(ctx context.Context, login string) (entities.Account, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT id, hash_pass FROM accounts WHERE login = $1;`
	var acc entities.Account
	err := a.pg.DB.QueryRow(ctx, q,
		login,
	).Scan(
		&acc.ID,
		&acc.HashPass,
	)
	if err != nil {
		return entities.Account{}, errors.NewErrInternal(err.Error())
	}

	acc.Login = login

	return acc, nil
}

func (a *AccountPG) GetAccountBalance(ctx context.Context, accountID string) (float64, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT points FROM accounts WHERE id = $1;`
	var balance float64
	err := a.pg.DB.QueryRow(ctx, q,
		accountID,
	).Scan(&balance)
	if err != nil {
		if storage.IsNotFound(err) {
			return 0, nil
		}

		return 0, errors.NewErrInternal(err.Error())
	}

	return balance, nil
}

func (a *AccountPG) UpdateAccountBalance(ctx context.Context, accountID string, delta float64) *errors.ErrorApp {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	tx, err := a.pg.DB.Begin(ctx)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !stdErr.Is(err, pgx.ErrTxClosed) {
			logger.Log.Error().Err(err).Msg("failed to rollback TX")
		}
	}()

	var balance float64
	q := `SELECT points FROM accounts WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, q, accountID).Scan(&balance)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	if balance+delta < 0 {
		return errors.NewErrInsufficientFunds()
	}

	q = `UPDATE accounts SET points = $2 WHERE id = $1`
	_, err = tx.Exec(ctx, q, accountID, balance+delta)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	return nil
}

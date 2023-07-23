package accunts

import (
	"context"

	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/storage"
	"gomarket/internal/storage/db"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, login string, hashPass string) (string, *errors.ErrorApp)
	GetAccountByLogin(ctx context.Context, login string) (entities.Account, *errors.ErrorApp)
	GetAccountBalance(ctx context.Context, accountID string) (int, *errors.ErrorApp)
	UpdateAccountBalance(ctx context.Context, accountID string, sum int) *errors.ErrorApp
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

func (a *AccountPG) GetAccountBalance(ctx context.Context, accountID string) (int, *errors.ErrorApp) {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	q := `SELECT points FROM accounts WHERE id = $1;`
	var balance int
	err := a.pg.DB.QueryRow(ctx, q,
		accountID,
	).Scan(&balance)
	if err != nil {
		return 0, errors.NewErrInternal(err.Error())
	}

	return balance, nil
}

func (a *AccountPG) UpdateAccountBalance(ctx context.Context, accountID string, sum int) *errors.ErrorApp {
	ctx, cancel := context.WithTimeout(ctx, storage.DefaultQueryTimeout)
	defer cancel()

	tx, err := a.pg.DB.Begin(ctx)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}
	defer tx.Rollback(ctx)

	var points int
	q := `SELECT points FROM accounts WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, q,
		accountID,
	).Scan(&points)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	if points+sum < 0 {
		return errors.NewErrInsufficientFunds()
	}

	q = `UPDATE accounts SET points = points + $2 WHERE id = $1 RETURNING points`
	err = tx.QueryRow(ctx, q,
		accountID,
		sum,
	).Scan(&points)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewErrInternal(err.Error())
	}

	return nil
}

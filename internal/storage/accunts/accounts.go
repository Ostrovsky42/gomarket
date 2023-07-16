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

package storage

import (
	"time"

	"github.com/jackc/pgconn"
)

const DefaultQueryTimeout = time.Second * 15

const UniqueViolationCode = "23505"

func IsUniqueViolation(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok && pgErr.Code == UniqueViolationCode {
		return true
	}
	return false
}

func IsNotFound(err error) bool {
	return err.Error() == "no rows in result set"
}

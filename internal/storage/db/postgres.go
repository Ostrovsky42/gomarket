package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"gomarket/internal/logger"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgresDB(dsn string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	pg := &Postgres{DB: pool}

	return pg, nil
}

func (p *Postgres) Close() {
	logger.Log.Info().Msg("CLOSE DB")
	p.DB.Close()
}

package db

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"gomarket/internal/logger"
)

const migrationPath = "./migrations"

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

	pg.upMigration()

	return pg, nil
}

func (p *Postgres) Close() {
	p.DB.Close()
	logger.Log.Info().Msg("db closed")
}

func (p *Postgres) upMigration() {
	db, err := p.DB.Acquire(context.Background())
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failde Acquire conn")

	}
	defer db.Release()

	sqlDB, err := sql.Open("postgres", db.Conn().Config().ConnString())
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed open db")
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create migrate driver")
	}

	fsrc, err := (&file.File{}).Open(migrationPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to open migrations directory")
	}

	m, err := migrate.NewWithInstance("file", fsrc, "postgres", driver)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create migrate instance")
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatal().Err(err).Msg("failed to apply migrations")
	}
}

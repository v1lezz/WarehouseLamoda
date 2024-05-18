package repository

import (
	"context"
	"errors"
	"fmt"
	"warehouse/internal/adapters/repository/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ( // errors
	ErrNotFound                   = errors.New("not found")
	ErrIsExist                    = errors.New("is exist")
	ErrIsNotExist                 = errors.New("is not exist")
	ErrFailedCheckGoodInWarehouse = errors.New("good in this warehouse does not exist")
	ErrReserve                    = errors.New("error reserve")
)

type PostgresConn struct {
	pool     *pgxpool.Pool
	isClosed bool
}

func NewPostgresConn(ctx context.Context, s string, migrationsPath string) (*PostgresConn, error) {
	cfg, err := pgxpool.ParseConfig(s)
	if err != nil {
		return nil, fmt.Errorf("error parse config from string: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("error create pool: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err = migrations.Up(pool, migrationsPath); err != nil {
		return nil, fmt.Errorf("error up migrations: %w", err)
	}
	return &PostgresConn{
		pool:     pool,
		isClosed: false,
	}, nil
}

func (pg *PostgresConn) Close() {
	if !pg.isClosed {
		pg.pool.Close()
	}
}

package repository

import (
	"context"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgUnitOfWorkFactory struct {
	pool *pgxpool.Pool
}

func NewPgUnitOfWorkFactory(pool *pgxpool.Pool) interfaces.UnitOfWorkFactory {
	return &PgUnitOfWorkFactory{pool: pool}
}

func (f *PgUnitOfWorkFactory) New(ctx context.Context) (interfaces.UnitOfWork, error) {
	tx, err := f.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return NewPgUnitOfWork(tx, ctx), nil
}

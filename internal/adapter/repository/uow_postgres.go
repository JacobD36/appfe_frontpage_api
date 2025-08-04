package repository

import (
	"context"
	"errors"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/jackc/pgx/v5"
)

type PgUnitOfWork struct {
	tx         pgx.Tx
	userRepo   interfaces.UserRepository
	committed  bool
	rolledBack bool
	ctx        context.Context
}

func NewPgUnitOfWork(tx pgx.Tx, ctx context.Context) *PgUnitOfWork {
	return &PgUnitOfWork{
		tx:       tx,
		userRepo: NewPgxUser(tx),
		ctx:      ctx,
	}
}

func (uow *PgUnitOfWork) Commit() error {
	if uow.committed {
		return errors.New(dto.ErrTransactionAlreadyCommitted)
	}
	if uow.rolledBack {
		return errors.New(dto.ErrTransactionAlreadyRolledBack)
	}

	err := uow.tx.Commit(uow.ctx)
	if err == nil {
		uow.committed = true
	}

	return err
}

func (uow *PgUnitOfWork) Rollback() error {
	if uow.committed {
		return errors.New(dto.ErrTransactionAlreadyCommitted)
	}
	if uow.rolledBack {
		return nil
	}

	err := uow.tx.Rollback(uow.ctx)
	if err == nil {
		uow.rolledBack = true
	}

	return err
}

func (uow *PgUnitOfWork) UserRepository() interfaces.UserRepository {
	return uow.userRepo
}

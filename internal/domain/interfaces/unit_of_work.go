package interfaces

import "context"

type UnitOfWork interface {
	Commit() error
	Rollback() error
	UserRepository() UserRepository
}

type UnitOfWorkFactory interface {
	New(ctx context.Context) (UnitOfWork, error)
}

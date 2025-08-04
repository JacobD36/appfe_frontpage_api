package usecase

import (
	"context"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
)

type MigrationService struct {
	uowFactory interfaces.UnitOfWorkFactory
}

func NewMigrationService(uowFactory interfaces.UnitOfWorkFactory) *MigrationService {
	return &MigrationService{uowFactory: uowFactory}
}

func (s *MigrationService) Migrate(ctx context.Context) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.UserRepository().Migrate(ctx); err != nil {
		return err
	}

	return uow.Commit()
}

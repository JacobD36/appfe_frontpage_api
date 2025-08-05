package usecase

import (
	"context"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	usecaseInterfaces "github.com/JacobD36/appfe_frontpage_api/internal/usecase/interfaces"
)

type MigrationService struct {
	uowFactory  interfaces.UnitOfWorkFactory
	userService usecaseInterfaces.UserService
}

func NewMigrationService(uowFactory interfaces.UnitOfWorkFactory, userService usecaseInterfaces.UserService) *MigrationService {
	return &MigrationService{
		uowFactory:  uowFactory,
		userService: userService,
	}
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

	if err := uow.Commit(); err != nil {
		return err
	}

	if err := s.userService.CreateInitialAdmin(ctx); err != nil {
		return err
	}

	return nil
}

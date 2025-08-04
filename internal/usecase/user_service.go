package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
	ui "github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/interfaces"
)

type userService struct {
	uowFactory ui.UnitOfWorkFactory
	hasher     ui.PasswordHasher
}

func NewUserService(uowFactory ui.UnitOfWorkFactory, h ui.PasswordHasher) interfaces.UserService {
	return &userService{
		uowFactory: uowFactory,
		hasher:     h,
	}
}

func (s *userService) Create(ctx context.Context, u *domain.User) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	validatedRole, err := domain.ValidateRole(u.Role)
	if err != nil {
		return err
	}
	u.Role = validatedRole

	u.CreatedAt = time.Now()
	u.Status = true
	u.Name = strings.ToUpper(strings.TrimSpace(u.Name))

	if u.Password != nil {
		hashed, err := s.hasher.Hash(*u.Password)
		if err != nil {
			return err
		}
		u.Password = &hashed
	}

	if err := uow.UserRepository().Create(ctx, u); err != nil {
		return err
	}
	return uow.Commit()
}

func (s *userService) UpdateByID(ctx context.Context, input ui.UpdateUserInput) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.UserRepository().UpdateByID(ctx, input); err != nil {
		return err
	}
	return uow.Commit()
}

func (s *userService) GetAll(ctx context.Context, pagination *domain.Pagination) (*domain.PaginatedResult[*domain.User], error) {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	users, total, err := uow.UserRepository().GetAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.Password = nil
	}

	if pagination == nil {
		pagination = domain.NewPagination(1, int(total), "")
	}

	result := domain.NewPaginatedResult(users, pagination, total)
	return result, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	user, err := uow.UserRepository().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	return uow.UserRepository().FindByEmail(ctx, email)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.UserRepository().Delete(ctx, id); err != nil {
		return err
	}
	return uow.Commit()
}

package interfaces

import (
	"context"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	UpdateByID(ctx context.Context, input interfaces.UpdateUserInput) error
	GetAll(ctx context.Context, pagination *domain.Pagination) (*domain.PaginatedResult[*domain.User], error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	CreateInitialAdmin(ctx context.Context) error
}

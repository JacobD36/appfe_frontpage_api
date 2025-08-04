package interfaces

import (
	"context"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
)

type UserRepository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, user *domain.User) error
	UpdateByID(ctx context.Context, input UpdateUserInput) error
	GetAll(ctx context.Context, pagination *domain.Pagination) ([]*domain.User, int64, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

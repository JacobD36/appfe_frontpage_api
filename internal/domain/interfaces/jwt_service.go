package interfaces

import "github.com/JacobD36/appfe_frontpage_api/internal/domain"

type JWTService interface {
	GenerateToken(user domain.User) (string, error)
	ValidateToken(tokenString string) (*domain.User, error)
}

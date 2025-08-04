package interfaces

import "github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"

type AuthService interface {
	Login(input dto.AuthLoginInput) (*dto.AuthLoginResponse, error)
	SignInWithToken(input dto.AuthTokenSignInInput) (*dto.AuthLoginResponse, error)
}

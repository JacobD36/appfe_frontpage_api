package usecase

import (
	"context"
	"errors"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
)

type AuthService struct {
	uowFactory     interfaces.UnitOfWorkFactory
	passwordHasher interfaces.PasswordHasher
	jwtService     interfaces.JWTService
}

func NewAuthService(
	uowFactory interfaces.UnitOfWorkFactory,
	passwordHasher interfaces.PasswordHasher,
	jwtService interfaces.JWTService,
) *AuthService {
	return &AuthService{
		uowFactory:     uowFactory,
		passwordHasher: passwordHasher,
		jwtService:     jwtService,
	}
}

func (s *AuthService) Login(input dto.AuthLoginInput) (*dto.AuthLoginResponse, error) {
	ctx := context.Background()

	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, errors.New(dto.ErrInternalServer)
	}
	defer uow.Rollback()

	userRepo := uow.UserRepository()

	user, err := userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if err.Error() == dto.ErrNoRowsFound {
			return nil, errors.New(dto.ErrInvalidCredentials)
		}
		return nil, errors.New(dto.ErrInternalServer)
	}

	if user.Password == nil {
		return nil, errors.New(dto.ErrInvalidCredentials)
	}

	if err := s.passwordHasher.Verify(*user.Password, input.Password); err != nil {
		return nil, errors.New(dto.ErrPasswordIncorrect)
	}

	if !user.EmailValidated {
		return nil, errors.New(dto.ErrEmailNotValidated)
	}

	if !user.Status {
		return nil, errors.New(dto.ErrAccountDisabled)
	}

	token, err := s.jwtService.GenerateToken(*user)
	if err != nil {
		return nil, errors.New(dto.ErrTokenGenerationFailed)
	}

	if err := uow.Commit(); err != nil {
		return nil, errors.New(dto.ErrInternalServer)
	}

	user.Password = nil

	response := &dto.AuthLoginResponse{
		Token: token,
		User:  *user,
	}

	return response, nil
}

func (s *AuthService) SignInWithToken(input dto.AuthTokenSignInInput) (*dto.AuthLoginResponse, error) {
	ctx := context.Background()

	user, err := s.jwtService.ValidateToken(input.Token)
	if err != nil {
		return nil, errors.New(dto.ErrInvalidToken)
	}

	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, errors.New(dto.ErrInternalServer)
	}
	defer uow.Rollback()

	userRepo := uow.UserRepository()

	fullUser, err := userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		if err.Error() == dto.ErrNoRowsFound {
			return nil, errors.New(dto.ErrUserNotFoundForToken)
		}
		return nil, errors.New(dto.ErrInternalServer)
	}

	if !fullUser.EmailValidated {
		return nil, errors.New(dto.ErrEmailNotValidated)
	}

	if !fullUser.Status {
		return nil, errors.New(dto.ErrAccountDisabled)
	}

	newToken, err := s.jwtService.GenerateToken(*fullUser)
	if err != nil {
		return nil, errors.New(dto.ErrTokenGenerationFailed)
	}

	if err := uow.Commit(); err != nil {
		return nil, errors.New(dto.ErrInternalServer)
	}

	fullUser.Password = nil

	response := &dto.AuthLoginResponse{
		Token: newToken,
		User:  *fullUser,
	}

	return response, nil
}

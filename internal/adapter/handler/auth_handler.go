package handler

import (
	"net/http"

	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/pkg/validator"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService interfaces.AuthService
}

func NewAuthHandler(authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input dto.AuthLoginInput
	if err := c.Bind(&input); err != nil {
		return Error(c, http.StatusBadRequest, dto.ErrInvalidInput)
	}

	if err := validator.Validate.Struct(input); err != nil {
		return Error(c, http.StatusBadRequest, dto.TranslateValidationErrors(err))
	}

	response, err := h.authService.Login(input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err.Error() {
		case dto.ErrInvalidCredentials:
			statusCode = http.StatusBadRequest
		case dto.ErrEmailNotValidated:
			statusCode = http.StatusBadRequest
		case dto.ErrAccountDisabled:
			statusCode = http.StatusBadRequest
		case dto.ErrPasswordIncorrect:
			statusCode = http.StatusBadRequest
		}

		return Error(c, statusCode, err.Error())
	}

	return SuccessLogin(c, http.StatusOK, dto.ErrLoginSuccess, response.User, response.Token)
}

func (h *AuthHandler) SignInWithToken(c echo.Context) error {
	var input dto.AuthTokenSignInInput
	if err := c.Bind(&input); err != nil {
		return Error(c, http.StatusBadRequest, dto.ErrInvalidInput)
	}

	if err := validator.Validate.Struct(input); err != nil {
		return Error(c, http.StatusBadRequest, dto.TranslateValidationErrors(err))
	}

	response, err := h.authService.SignInWithToken(input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err.Error() {
		case dto.ErrInvalidToken:
			statusCode = http.StatusUnauthorized
		case dto.ErrUserNotFoundForToken:
			statusCode = http.StatusBadRequest
		case dto.ErrEmailNotValidated:
			statusCode = http.StatusBadRequest
		case dto.ErrAccountDisabled:
			statusCode = http.StatusBadRequest
		}

		return Error(c, statusCode, err.Error())
	}

	return SuccessLogin(c, http.StatusOK, dto.ErrTokenSignInSuccess, response.User, response.Token)
}

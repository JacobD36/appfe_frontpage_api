package handler

import (
	"net/http"

	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/pkg/logger"
	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data"`
}

type APILoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data"`
	Token   string `json:"token"`
}

func Success(c echo.Context, code int, message string, data any) error {
	ctx := c.Request().Context()
	logger.Info(ctx, dto.MsgSuccessfulResponse,
		logger.Int("status_code", code),
		logger.String("message", message),
	)

	return c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Status:  http.StatusText(code),
		Data:    data,
	})
}

func SuccessLogin(c echo.Context, code int, message string, data any, token string) error {
	ctx := c.Request().Context()
	logger.Info(ctx, dto.MsgSuccessfulLoginResponse,
		logger.Int("status_code", code),
		logger.String("message", message),
		logger.String("token_length", string(rune(len(token)))),
	)

	return c.JSON(code, APILoginResponse{
		Code:    code,
		Message: message,
		Status:  http.StatusText(code),
		Data:    data,
		Token:   token,
	})
}

func Error(c echo.Context, code int, message string) error {
	ctx := c.Request().Context()
	logger.Warn(ctx, dto.MsgErrorResponse,
		logger.Int("status_code", code),
		logger.String("message", message),
		logger.String("path", c.Request().URL.Path),
		logger.String("method", c.Request().Method),
	)

	return c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Status:  http.StatusText(code),
		Data:    nil,
	})
}

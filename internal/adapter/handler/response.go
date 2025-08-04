package handler

import (
	"net/http"

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
	return c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Status:  http.StatusText(code),
		Data:    data,
	})
}

func SuccessLogin(c echo.Context, code int, message string, data any, token string) error {
	return c.JSON(code, APILoginResponse{
		Code:    code,
		Message: message,
		Status:  http.StatusText(code),
		Data:    data,
		Token:   token,
	})
}

func Error(c echo.Context, code int, message string) error {
	return c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Status:  http.StatusText(code),
		Data:    nil,
	})
}

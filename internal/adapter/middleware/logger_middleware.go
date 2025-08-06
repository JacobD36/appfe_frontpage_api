package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/pkg/logger"
	"github.com/labstack/echo/v4"
)

// contextKey es un tipo personalizado para evitar colisiones en el contexto
type contextKey string

const (
	// RequestIDKey es la clave para almacenar el request ID en el contexto
	RequestIDKey contextKey = "request_id"
)

type LoggerConfig struct {
	SkipURIs []string

	// LogErrorResponseBody indica si se debe loggear el cuerpo de las respuestas de error
	LogErrorResponseBody bool
}

// DefaultLoggerConfig es la configuración por defecto
var DefaultLoggerConfig = LoggerConfig{
	SkipURIs:             []string{"/health", "/metrics"},
	LogErrorResponseBody: false,
}

// Logger crea un middleware de logging personalizado
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithConfig crea un middleware de logging con configuración personalizada
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Verificar si la URI debe ser omitida
			for _, skipURI := range config.SkipURIs {
				if c.Request().URL.Path == skipURI {
					return next(c)
				}
			}

			// Generar request ID si no existe
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = generateRequestID()
				c.Response().Header().Set(echo.HeaderXRequestID, requestID)
			}

			// Agregar request ID al contexto para que esté disponible en los logs
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))

			// Log del inicio de la request
			logger.Info(ctx, dto.MsgRequestStarted,
				logger.String("method", c.Request().Method),
				logger.String("path", c.Request().URL.Path),
				logger.String("remote_addr", c.RealIP()),
				logger.String("user_agent", c.Request().UserAgent()),
				logger.String("request_id", requestID),
			)

			// Ejecutar el siguiente handler
			err := next(c)

			// Calcular duración
			duration := time.Since(start)

			// Obtener status code
			statusCode := c.Response().Status
			if err != nil {
				// Si hay error, Echo manejará el status code
				if he, ok := err.(*echo.HTTPError); ok {
					statusCode = he.Code
				} else {
					statusCode = 500
				}
			}

			// Log de la respuesta
			logger.GetLogger().LogHTTPRequest(c, duration, statusCode)

			// Log adicional para errores
			if err != nil {
				logger.GetLogger().Error(ctx, dto.MsgRequestError,
					logger.String("method", c.Request().Method),
					logger.String("path", c.Request().URL.Path),
					logger.Int("status_code", statusCode),
					logger.Duration("duration", duration),
					logger.Error("error", err),
				)
			}

			return err
		}
	}
}

// generateRequestID genera un ID único para la request
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

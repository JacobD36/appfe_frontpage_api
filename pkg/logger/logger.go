package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/labstack/echo/v4"
)

// LogLevel representa los niveles de logging disponibles
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// Logger es una interfaz para el sistema de logging
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)

	// Métodos específicos para el dominio de la aplicación
	LogHTTPRequest(c echo.Context, duration time.Duration, statusCode int)
	LogDatabaseOperation(ctx context.Context, operation string, table string, duration time.Duration, err error)
	LogAuthentication(ctx context.Context, email string, success bool, reason string)
	LogBusinessOperation(ctx context.Context, operation string, userID string, details map[string]interface{})
}

// Field representa un campo adicional para el log
type Field struct {
	Key   string
	Value any
}

// appLogger implementa la interfaz Logger usando slog
type appLogger struct {
	logger *slog.Logger
	level  slog.Level
}

// New crea una nueva instancia del logger
func New(level LogLevel) Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Personalizar el formato del timestamp
			if a.Key == slog.TimeKey {
				return slog.String("timestamp", a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &appLogger{
		logger: logger,
		level:  slogLevel,
	}
}

// Helper functions para crear campos de log
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func Error(key string, value error) Field {
	if value == nil {
		return Field{Key: key, Value: nil}
	}
	return Field{Key: key, Value: value.Error()}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Implementación de los métodos de logging
func (l *appLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, slog.LevelDebug, msg, fields...)
}

func (l *appLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, slog.LevelInfo, msg, fields...)
}

func (l *appLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, slog.LevelWarn, msg, fields...)
}

func (l *appLogger) Error(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, slog.LevelError, msg, fields...)
}

func (l *appLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, slog.LevelError, msg, fields...)
	os.Exit(1)
}

// logWithContext agrega contexto adicional al log
func (l *appLogger) logWithContext(ctx context.Context, level slog.Level, msg string, fields ...Field) {
	if !l.logger.Enabled(ctx, level) {
		return
	}

	attrs := make([]slog.Attr, 0, len(fields)+2)

	// Agregar información del contexto si está disponible
	if ctx != nil {
		if requestID := getRequestID(ctx); requestID != "" {
			attrs = append(attrs, slog.String("request_id", requestID))
		}
		if userID := getUserID(ctx); userID != "" {
			attrs = append(attrs, slog.String("user_id", userID))
		}
	}

	// Agregar campos adicionales
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}

	l.logger.LogAttrs(ctx, level, msg, attrs...)
}

// Métodos específicos del dominio
func (l *appLogger) LogHTTPRequest(c echo.Context, duration time.Duration, statusCode int) {
	ctx := c.Request().Context()

	level := slog.LevelInfo
	if statusCode >= 400 {
		level = slog.LevelWarn
	}
	if statusCode >= 500 {
		level = slog.LevelError
	}

	l.logWithContext(ctx, level, dto.MsgHTTPRequest,
		String("method", c.Request().Method),
		String("path", c.Request().URL.Path),
		String("remote_addr", c.RealIP()),
		String("user_agent", c.Request().UserAgent()),
		Int("status_code", statusCode),
		Duration("duration", duration),
		Any("query_params", c.QueryParams()),
	)
}

func (l *appLogger) LogDatabaseOperation(ctx context.Context, operation string, table string, duration time.Duration, err error) {
	if err != nil {
		l.logWithContext(ctx, slog.LevelError, dto.MsgDatabaseOpFailed,
			String("operation", operation),
			String("table", table),
			Duration("duration", duration),
			Error("error", err),
		)
	} else {
		l.logWithContext(ctx, slog.LevelDebug, dto.MsgDatabaseOpSuccess,
			String("operation", operation),
			String("table", table),
			Duration("duration", duration),
		)
	}
}

func (l *appLogger) LogAuthentication(ctx context.Context, email string, success bool, reason string) {
	level := slog.LevelInfo
	msg := dto.MsgAuthenticationSuccess

	if !success {
		level = slog.LevelWarn
		msg = dto.MsgAuthenticationFailed
	}

	l.logWithContext(ctx, level, msg,
		String("email", email),
		String("reason", reason),
		Any("success", success),
	)
}

func (l *appLogger) LogBusinessOperation(ctx context.Context, operation string, userID string, details map[string]interface{}) {
	fields := []Field{
		String("operation", operation),
		String("user_id", userID),
	}

	for key, value := range details {
		fields = append(fields, Any(key, value))
	}

	l.logWithContext(ctx, slog.LevelInfo, dto.MsgBusinessOperation, fields...)
}

// Helper functions para extraer información del contexto
func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Intentar obtener del contexto de Echo
	if ec, ok := ctx.Value("echo").(*echo.Context); ok {
		return (*ec).Response().Header().Get(echo.HeaderXRequestID)
	}

	// Intentar obtener directamente del contexto
	if reqID, ok := ctx.Value("request_id").(string); ok {
		return reqID
	}

	return ""
}

func getUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Intentar obtener del contexto de Echo
	if ec, ok := ctx.Value("echo").(*echo.Context); ok {
		if userID := (*ec).Get("user_id"); userID != nil {
			if id, ok := userID.(string); ok {
				return id
			}
		}
	}

	// Intentar obtener directamente del contexto
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}

	return ""
}

// Instancia global del logger
var defaultLogger Logger

// Init inicializa el logger global
func Init(level LogLevel) {
	defaultLogger = New(level)
}

// GetLogger retorna la instancia global del logger
func GetLogger() Logger {
	if defaultLogger == nil {
		defaultLogger = New(LevelInfo)
	}
	return defaultLogger
}

// Funciones de conveniencia para usar el logger global
func Debug(ctx context.Context, msg string, fields ...Field) {
	GetLogger().Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	GetLogger().Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	GetLogger().Warn(ctx, msg, fields...)
}

func LogError(ctx context.Context, msg string, fields ...Field) {
	GetLogger().Error(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...Field) {
	GetLogger().Fatal(ctx, msg, fields...)
}

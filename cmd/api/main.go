package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/messaging"
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/router"
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/security"
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/storage"
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/template"
	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/pkg/logger"
	"github.com/joho/godotenv"
)

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	logger.Init(logger.LogLevel(logLevel))

	// Intentar cargar .env desde diferentes ubicaciones
	envPaths := []string{".env", "../.env", "../../.env"}
	var loadErr error
	for _, path := range envPaths {
		if loadErr = godotenv.Load(path); loadErr == nil {
			break
		}
	}
	if loadErr != nil {
		logger.Warn(context.Background(), dto.ErrLoadingEnvFile, logger.Error("error", loadErr))
	}
}

func initKeys() {
	ctx := context.Background()
	privateKeyPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	if privateKeyPath == "" || publicKeyPath == "" {
		logger.Fatal(ctx, dto.ErrRSAKeysNotSet)
	}

	logger.Info(ctx, dto.MsgLoadingRSAKeys,
		logger.String("private_key_path", privateKeyPath),
		logger.String("public_key_path", publicKeyPath),
	)

	err := security.LoadFiles(privateKeyPath, publicKeyPath)
	if err != nil {
		logger.Fatal(ctx, dto.ErrFailedLoadRSAKeys, logger.Error("error", err))
	}

	logger.Info(ctx, dto.MsgRSAKeysLoadedSuccess)
}

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	logger.Info(ctx, dto.MsgStartingAPIServer, logger.String("port", port))

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	driver := storage.Postgres
	logger.Info(ctx, dto.MsgInitializingDBConnection, logger.String("driver", string(driver)))
	storage.New(driver)

	initKeys()

	uowFactory, err := storage.UoWFactory(driver)
	if err != nil {
		logger.Fatal(ctx, dto.ErrUnitOfWorkFactory, logger.Error("error", err))
	}

	hasher := security.NewBcryptHasher(12)

	// Inicializar servicio de mensajería
	var messagingService interfaces.MessagingService
	if os.Getenv("BREVO_API_KEY") != "" {
		emailService, err := messaging.NewBrevoEmailServiceFromEnv(logger.GetLogger())
		if err != nil {
			logger.Warn(ctx, dto.ErrMessagingServiceInitFailed, logger.Error("error", err))
		} else {
			messagingService = usecase.NewMessagingService(emailService, logger.GetLogger())
			logger.Info(ctx, dto.MsgMessagingServiceInitSuccess)
		}
	} else {
		logger.Info(ctx, dto.MsgMessagingServiceDisabled)
	}

	// Inicializar servicio de templates
	templateService := template.NewHTMLTemplateService()
	logger.Info(ctx, dto.MsgTemplateServiceInitialized)

	userService := usecase.NewUserService(uowFactory, hasher, messagingService, templateService)

	logger.Info(ctx, dto.MsgRunningDBMigrations)
	migrationService := usecase.NewMigrationService(uowFactory, userService)
	if err := migrationService.Migrate(context.Background()); err != nil {
		logger.Fatal(ctx, dto.ErrMigrationFailed, logger.Error("error", err))
	}
	logger.Info(ctx, dto.MsgDBMigrationsCompleted)

	jwtService, err := security.NewJWTService()
	if err != nil {
		logger.Fatal(ctx, dto.ErrTokenGenerationFailed, logger.Error("error", err))
	}

	authService := usecase.NewAuthService(uowFactory, hasher, jwtService)

	r := router.New(userService, authService, jwtService)

	logger.Info(ctx, dto.MsgServicesInitialized)

	go func() {
		logger.Info(ctx, dto.MsgStartingHTTPServer, logger.String("address", port))
		if err := r.Start(port); err != nil {
			logger.LogError(ctx, dto.ErrServerError, logger.Error("error", err))
			stop()
		}
	}()

	<-signalCtx.Done()
	logger.Info(ctx, dto.MsgShutdownSignalReceived)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.GetEchoInstance().Shutdown(shutdownCtx); err != nil {
		logger.Fatal(ctx, dto.ErrForcedShutdown, logger.Error("error", err))
	}

	logger.Info(ctx, dto.MsgServerStoppedGracefully)
}

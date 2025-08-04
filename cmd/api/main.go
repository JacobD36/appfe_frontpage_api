package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/router"
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/security"
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/storage"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println(dto.ErrLoadingEnvFile)
	}
}

func initKeys() {
	privateKeyPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	if privateKeyPath == "" || publicKeyPath == "" {
		log.Fatal(dto.ErrRSAKeysNotSet)
	}
	err := security.LoadFiles(privateKeyPath, publicKeyPath)
	if err != nil {
		log.Fatalf(dto.ErrFailedLoadRSAKeys, err)
	}
}

func main() {
	port := os.Getenv("PORT")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	driver := storage.Postgres
	storage.New(driver)
	initKeys()

	uowFactory, err := storage.UoWFactory(driver)
	if err != nil {
		log.Fatal(dto.ErrUnitOfWorkFactory, ": ", err)
	}

	hasher := security.NewBcryptHasher(12)
	userService := usecase.NewUserService(uowFactory, hasher)

	migrationService := usecase.NewMigrationService(uowFactory, userService)
	if err := migrationService.Migrate(context.Background()); err != nil {
		log.Fatalf(dto.ErrMigrationFailed, err)
	}

	jwtService, err := security.NewJWTService()
	if err != nil {
		log.Fatal(dto.ErrTokenGenerationFailed, ": ", err)
	}

	authService := usecase.NewAuthService(uowFactory, hasher, jwtService)

	r := router.New(userService, authService, jwtService)

	go func() {
		if err := r.Start(port); err != nil {
			log.Printf(dto.ErrServerError, err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println(dto.MsgShuttingDownServer)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()

	if err := r.GetEchoInstance().Shutdown(shutdownCtx); err != nil {
		log.Fatalf(dto.ErrForcedShutdown, err)
	}

	log.Println(dto.MsgServerStoppedGracefully)
}

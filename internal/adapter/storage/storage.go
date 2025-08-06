package storage

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/repository"
	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

type Driver string

const (
	Postgres Driver = "POSTGRES"
)

func New(d Driver) {
	ctx := context.Background()
	switch d {
	case Postgres:
		logger.Info(ctx, dto.MsgInitializingPostgresConn)
		newPostgresConn()
	default:
		logger.Fatal(ctx, dto.ErrUnsupportedDriver, logger.String("driver", string(d)))
	}
}

func newPostgresConn() {
	once.Do(func() {
		ctx := context.Background()
		dsn := os.Getenv("POSTGRES_DATABASE_URL")
		if dsn == "" {
			logger.Fatal(ctx, dto.ErrDatabaseURLNotSet)
		}

		logger.Debug(ctx, dto.MsgParsingDBConfig)
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			logger.Fatal(ctx, dto.ErrFailedParseConfig, logger.Error("error", err))
		}

		logger.Debug(ctx, dto.MsgCreatingDBPool)
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			logger.Fatal(ctx, dto.ErrUnableCreatePool, logger.Error("error", err))
		}

		logger.Info(ctx, dto.MsgConnectedToDBSuccess,
			logger.String("host", config.ConnConfig.Host),
			logger.Int("port", int(config.ConnConfig.Port)),
			logger.String("database", config.ConnConfig.Database),
		)
	})
}

func Pool() *pgxpool.Pool {
	return pool
}

func UoWFactory(driver Driver) (interfaces.UnitOfWorkFactory, error) {
	switch driver {
	case Postgres:
		return repository.NewPgUnitOfWorkFactory(pool), nil
	default:
		return nil, fmt.Errorf(dto.ErrDriverNotImplemented, driver)
	}
}

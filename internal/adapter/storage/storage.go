package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/repository"
	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
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
	switch d {
	case Postgres:
		newPostgresConn()
	default:
		log.Fatalf(dto.ErrUnsupportedDriver, d)
	}
}

func newPostgresConn() {
	once.Do(func() {
		dsn := os.Getenv("POSTGRES_DATABASE_URL")
		if dsn == "" {
			log.Fatal(dto.ErrDatabaseURLNotSet)
		}

		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			log.Fatalf(dto.ErrFailedParseConfig, err)
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatalf(dto.ErrUnableCreatePool, err)
		}

		fmt.Println(dto.MsgConnectedToDatabase)
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

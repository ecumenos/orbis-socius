package datastore

import (
	"context"

	"github.com/ecumenos/orbis-socius/cmd/api/configuration"
	"github.com/ecumenos/orbis-socius/internal/postgres"
	"github.com/jackc/pgx/v4"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func(lc fx.Lifecycle, cfg *configuration.Config) (Driver, error) {
		driver, err := postgres.New(context.Background(), cfg.APIDataStore.URL)
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return driver.Ping(ctx)
			},
			OnStop: func(context.Context) error {
				driver.Close()
				return nil
			},
		})

		return driver, nil
	}),
)

//go:generate mockery --name=Driver

type Driver interface {
	Ping(ctx context.Context) error
	Close()
	CountRows(ctx context.Context, query string, args ...interface{}) (int, error)
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) error
	QueryRow(ctx context.Context, query string, args ...interface{}) (pgx.Row, error)
	QueryRows(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

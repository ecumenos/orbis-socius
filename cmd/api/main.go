package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/carlmjohnson/versioninfo"
	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxlogger"
	"github.com/ecumenos/fxecumenos/fxpostgres"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/fxecumenos/zerodowntime"
	"github.com/ecumenos/orbis-socius/cmd/api/accounts"
	"github.com/ecumenos/orbis-socius/cmd/api/configuration"
	"github.com/ecumenos/orbis-socius/cmd/api/httpserver"
	"github.com/ecumenos/orbis-socius/internal/postgres"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

var Version = fxecumenos.NewVersion("0.0.1")

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("exiting", "err", err)
		os.Exit(-1)
	}
}

func run(args []string) error {
	app := cli.App{
		Name:    "api",
		Usage:   "serving API",
		Version: versioninfo.Short(),
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "prod",
			Value:   false,
			EnvVars: []string{"PROD"},
		},
	}

	app.Commands = []*cli.Command{
		runAppCmd,
		migrateUpCmd,
		migrateDownCmd,
	}

	return app.Run(args)
}

var runAppCmd = &cli.Command{
	Name:  "run-api-server",
	Usage: "run API HTTP server",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Supply(fxecumenos.ServiceName("api")),
			fx.Supply(fxecumenos.Version(Version)),
			fxlogger.Module,
			fxrf.Module,
			fxpostgres.Module,
			configuration.Module,
			httpserver.Module,
			accounts.Module,
		)

		return zerodowntime.HandleApp(app)
	},
}

var migrateUpCmd = &cli.Command{
	Name:  "migrate-up",
	Usage: "run migrations up",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Supply(fxecumenos.ServiceName("api")),
			fx.Supply(fxecumenos.Version(Version)),
			fxlogger.Module,
			configuration.Module,
			fx.Invoke(func(cfg *configuration.Config, logger *zap.Logger, shutdowner fx.Shutdowner) error {
				fn := postgres.NewMigrateUpFunc()
				if !cctx.Bool("prod") {
					logger.Info("runnning migrate up",
						zap.String("db_url", cfg.APIDataStore.URL),
						zap.String("source_path", cfg.APIDataStore.MigrationsPath))
				}
				return fn(cfg.APIDataStore.MigrationsPath, cfg.APIDataStore.URL+"?sslmode=disable", logger, shutdowner)
			}),
		)

		return zerodowntime.HandleApp(app)
	},
}

var migrateDownCmd = &cli.Command{
	Name:  "migrate-down",
	Usage: "run migrations down",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Supply(fxecumenos.ServiceName("api")),
			fx.Supply(fxecumenos.Version(Version)),
			fxlogger.Module,
			configuration.Module,
			fx.Invoke(func(cfg *configuration.Config, logger *zap.Logger, shutdowner fx.Shutdowner) error {
				fn := postgres.NewMigrateDownFunc()
				if !cctx.Bool("prod") {
					logger.Info("runnning migrate down",
						zap.String("db_url", cfg.APIDataStore.URL),
						zap.String("source_path", cfg.APIDataStore.MigrationsPath))
				}
				return fn(cfg.APIDataStore.MigrationsPath, cfg.APIDataStore.URL+"?sslmode=disable", logger, shutdowner)
			}),
		)

		return zerodowntime.HandleApp(app)
	},
}

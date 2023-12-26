package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/carlmjohnson/versioninfo"
	"github.com/ecumenos/orbis-socius/cmd/adminmanager/configuration"
	"github.com/ecumenos/orbis-socius/cmd/adminmanager/httpserver"
	"github.com/ecumenos/orbis-socius/pkg/ecumenosfx"
	"github.com/ecumenos/orbis-socius/pkg/logger"
	"github.com/ecumenos/orbis-socius/pkg/zerodowntime"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("exiting", "err", err)
		os.Exit(-1)
	}
}

func run(args []string) error {
	app := cli.App{
		Name:    "admin-manager",
		Usage:   "serving administration management API",
		Version: versioninfo.Short(),
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "prod",
			Value:   false,
			EnvVars: []string{"PROD"},
		},
	}

	app.Action = func(cctx *cli.Context) error {
		app := fx.New(
			fx.Supply(ecumenosfx.ServiceName("admin-manager")),
			fx.Provide(
				func(lc fx.Lifecycle, sn ecumenosfx.ServiceName) (*zap.Logger, error) {
					return logger.NewZapLogger(sn, cctx.Bool("prod"), lc)
				},
				logger.ZapSugared,
			),
			fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			}),
			configuration.Module,
			httpserver.Module,
		)

		return zerodowntime.HandleApp(app)
	}

	return app.Run(args)
}

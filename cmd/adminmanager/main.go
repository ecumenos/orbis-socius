package main

import (
	"os"
	"strings"

	"go.uber.org/fx"

	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxlogger"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/fxecumenos/zerodowntime"
	"github.com/ecumenos/orbis-socius/cmd/adminmanager/configuration"
	"github.com/ecumenos/orbis-socius/cmd/adminmanager/httpserver"
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
		Name:    "admin-manager",
		Usage:   "serving administration management API",
		Version: strings.Join(Version.Build, "."),
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
			fx.Supply(fxecumenos.ServiceName("admin-manager")),
			fx.Supply(fxecumenos.Version(Version)),
			fxlogger.Module,
			fxrf.Module,
			configuration.Module,
			httpserver.Module,
		)

		return zerodowntime.HandleApp(app)
	}

	return app.Run(args)
}

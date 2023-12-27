package configuration

import (
	"github.com/ecumenos/fxecumenos/fxlogger"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/go-toolkit/envutils"
	"github.com/ecumenos/orbis-socius/cmd/adminmanager/httpserver"
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func() (*Config, error) {
		var cfg Config

		if envutils.GetEnvBoolWithDefault("ADMIN_MANAGER_LOCAL", false) {
			if err := godotenv.Load(); err != nil {
				return nil, err
			}
		}

		if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&cfg, "cmd/adminmanager/configuration/default.json"); err != nil {
			return nil, err
		}

		return &cfg, nil
	},

		func(cfg *Config) *httpserver.Config { return cfg.AdminManagerHTTP },
		func(cfg *Config) *fxlogger.Config { return cfg.AdminManagerLogger },
		func(cfg *Config) *fxrf.Config { return cfg.AdminManagerResponseFactory },
	),
)

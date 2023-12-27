package configuration

import (
	"github.com/ecumenos/fxecumenos/fxlogger"
	"github.com/ecumenos/fxecumenos/fxpostgres"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/go-toolkit/envutils"
	"github.com/ecumenos/orbis-socius/cmd/api/httpserver"
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (*Config, error) {
			var cfg Config

			if envutils.GetEnvBoolWithDefault("API_LOCAL", false) {
				if err := godotenv.Load(); err != nil {
					return nil, err
				}
			}

			if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&cfg, "cmd/api/configuration/default.json"); err != nil {
				return nil, err
			}

			return &cfg, nil
		},
		func(cfg *Config) *httpserver.Config { return cfg.APIHTTP },
		func(cfg *Config) *fxpostgres.Config { return cfg.APIDataStore },
		func(cfg *Config) *fxlogger.Config { return cfg.APILogger },
		func(cfg *Config) *fxrf.Config { return cfg.APIResponseFactory },
	),
)

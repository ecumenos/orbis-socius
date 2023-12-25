package configuration

import (
	"errors"
	"os"
	"strconv"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func() (*Config, error) {
		var cfg Config

		if getenvBoolWithDefault("API_LOCAL", false) {
			if err := godotenv.Load(); err != nil {
				return nil, err
			}
		}

		if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&cfg, "cmd/api/configuration/default.json"); err != nil {
			return nil, err
		}

		return &cfg, nil
	}),
)

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, errors.New("getenv: environment variable empty")
	}
	return v, nil
}

func getenvBool(key string) (bool, error) {
	s, err := getenvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func getenvBoolWithDefault(key string, def bool) bool {
	s, err := getenvStr(key)
	if err != nil {
		return def
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return v
}

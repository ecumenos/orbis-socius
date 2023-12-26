package configuration

import "time"

type Config struct {
	APIProd      bool
	APILocal     bool `default:"true"`
	APIHTTP      HTTPConfig
	APIDataStore DataStoreConfig
}

type HTTPConfig struct {
	Addr           string        `default:":9090"`
	HandlerTimeout time.Duration `default:"30s"`
	ReadTimeout    time.Duration `default:"15s"`
	WriteTimeout   time.Duration `default:"15s"`
	IdleTimeout    time.Duration `default:"15s"`
}

type DataStoreConfig struct {
	URL            string `json:"url" default:"postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_orbis_socius_apidb"`
	MigrationsPath string `json:"migrationsPath" default:"file://cmd/api/repo/migrations"`
}

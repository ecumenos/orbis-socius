package configuration

import "time"

type Config struct {
	APIProd  bool
	APILocal bool `default:"true"`
	APIHTTP  HTTPConfig
}

type HTTPConfig struct {
	Addr           string        `default:":9090"`
	HandlerTimeout time.Duration `default:"30s"`
	ReadTimeout    time.Duration `default:"15s"`
	WriteTimeout   time.Duration `default:"15s"`
	IdleTimeout    time.Duration `default:"15s"`
}

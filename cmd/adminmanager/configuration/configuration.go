package configuration

import "time"

type Config struct {
	AdminManagerProd  bool
	AdminManagerLocal bool `default:"true"`
	AdminManagerHTTP  HTTPConfig
}

type HTTPConfig struct {
	Addr           string        `default:":9091"`
	HandlerTimeout time.Duration `default:"30s"`
	ReadTimeout    time.Duration `default:"15s"`
	WriteTimeout   time.Duration `default:"15s"`
	IdleTimeout    time.Duration `default:"15s"`
}

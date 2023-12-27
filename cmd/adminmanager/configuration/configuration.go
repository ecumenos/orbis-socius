package configuration

import (
	"github.com/ecumenos/fxecumenos/fxlogger"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/orbis-socius/cmd/adminmanager/httpserver"
)

type Config struct {
	AdminManagerLocal           bool `default:"true"`
	AdminManagerHTTP            *httpserver.Config
	AdminManagerLogger          *fxlogger.Config
	AdminManagerResponseFactory *fxrf.Config
}

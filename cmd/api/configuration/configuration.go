package configuration

import (
	"github.com/ecumenos/fxecumenos/fxlogger"
	"github.com/ecumenos/fxecumenos/fxpostgres"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/orbis-socius/cmd/api/httpserver"
)

type Config struct {
	APILocal           bool `default:"true"`
	APILogger          *fxlogger.Config
	APIHTTP            *httpserver.Config
	APIDataStore       *fxpostgres.Config
	APIResponseFactory *fxrf.Config
}

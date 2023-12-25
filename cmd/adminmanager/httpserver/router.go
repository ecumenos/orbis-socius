package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type routesParams struct {
	fx.In

	Handlers Handlers
	Logger   *zap.Logger
}

func NewRouter(params routesParams) *mux.Router {
	router := mux.NewRouter()
	enrichContext := NewEnrichContextMiddleware(params.Logger)
	logRequest := NewLoggerMiddleware(params.Logger)
	recovery := NewRecoverMiddleware(params.Logger)

	router.Use(mux.MiddlewareFunc(enrichContext))
	router.HandleFunc("/ping", params.Handlers.Ping).Methods(http.MethodGet)
	router.HandleFunc("/info", params.Handlers.Info).Methods(http.MethodGet)
	router.Use(mux.MiddlewareFunc(logRequest))
	router.Use(mux.MiddlewareFunc(recovery))

	return router
}

package httpserver

import (
	"net/http"

	"github.com/ecumenos/orbis-socius/pkg/ecumenosfx"
	"github.com/ecumenos/orbis-socius/pkg/toolkit/httputils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handlers interface {
	Ping(rw http.ResponseWriter, r *http.Request)
	Info(rw http.ResponseWriter, r *http.Request)
}

type HandlersImpl struct {
	Logger *zap.Logger
	Name   ecumenosfx.ServiceName
}

type handlersParams struct {
	fx.In

	Logger *zap.Logger
	Name   ecumenosfx.ServiceName
}

func NewHandlers(params handlersParams) Handlers {
	return &HandlersImpl{
		Logger: params.Logger,
		Name:   params.Name,
	}
}

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (h *HandlersImpl) Ping(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := httputils.NewWriter(h.Logger, rw)
	writer.WriteSuccess(ctx, &GetPingRespData{Ok: true})
}

type GetInfoRespData struct {
	Name string `json:"name"`
}

func (h *HandlersImpl) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := httputils.NewWriter(h.Logger, rw)
	writer.WriteSuccess(ctx, &GetInfoRespData{
		Name: string(h.Name),
	})
}

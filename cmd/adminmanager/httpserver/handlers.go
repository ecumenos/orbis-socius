package httpserver

import (
	"net/http"

	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxrf"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:generate mockery --name=Handlers

type Handlers interface {
	Ping(rw http.ResponseWriter, r *http.Request)
	Info(rw http.ResponseWriter, r *http.Request)
}

type HandlersImpl struct {
	Logger          *zap.Logger
	Name            fxecumenos.ServiceName
	ResponseFactory fxrf.Factory
}

type handlersParams struct {
	fx.In

	Logger          *zap.Logger
	Name            fxecumenos.ServiceName
	ResponseFactory fxrf.Factory
}

func NewHandlers(params handlersParams) Handlers {
	return &HandlersImpl{
		Logger:          params.Logger,
		Name:            params.Name,
		ResponseFactory: params.ResponseFactory,
	}
}

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (h *HandlersImpl) Ping(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.ResponseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetPingRespData{Ok: true})
}

type GetInfoRespData struct {
	Name string `json:"name"`
}

func (h *HandlersImpl) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.ResponseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetInfoRespData{
		Name: string(h.Name),
	})
}

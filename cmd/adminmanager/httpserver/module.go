package httpserver

import (
	"context"
	"net/http"

	"github.com/ecumenos/orbis-socius/cmd/adminmanager/configuration"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(NewHandlers, NewRouter, NewServer),
	fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner, s *Server) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					_ = s.Start(shutdowner)
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.Stop(ctx)
			},
		})
	}),
)

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(cfg *configuration.Config, logger *zap.Logger, router *mux.Router) *Server {
	return &Server{
		server: &http.Server{
			Addr:         cfg.AdminManagerHTTP.Addr,
			WriteTimeout: cfg.AdminManagerHTTP.WriteTimeout,
			ReadTimeout:  cfg.AdminManagerHTTP.ReadTimeout,
			IdleTimeout:  cfg.AdminManagerHTTP.IdleTimeout,
			Handler:      http.TimeoutHandler(router, cfg.AdminManagerHTTP.HandlerTimeout, "something went wrong"),
		},
		logger: logger,
	}

}

func (s *Server) Start(shutdowner fx.Shutdowner) error {
	s.logger.Info("http server is starting...")
	shutdownStatus := s.server.ListenAndServe()
	s.logger.Info("http server shutdown status", zap.Any("status", shutdownStatus))
	if err := shutdowner.Shutdown(); err != nil {
		s.logger.Error("shutdown http server error", zap.Error(err))
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("http server is shutting down...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("shutting down http server error", zap.Error(err))
		return err
	}
	return nil
}

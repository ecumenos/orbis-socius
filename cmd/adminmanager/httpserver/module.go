package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	Addr           string        `default:":9091"`
	HandlerTimeout time.Duration `default:"30s"`
	ReadTimeout    time.Duration `default:"15s"`
	WriteTimeout   time.Duration `default:"15s"`
	IdleTimeout    time.Duration `default:"15s"`
}

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

func NewServer(cfg *Config, logger *zap.Logger, router *mux.Router) *Server {
	return &Server{
		server: &http.Server{
			Addr:         cfg.Addr,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			Handler:      http.TimeoutHandler(router, cfg.HandlerTimeout, "something went wrong"),
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

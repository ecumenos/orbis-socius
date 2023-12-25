package logger

import (
	"time"

	"github.com/ecumenos/orbis-socius/pkg/toolkit/timeutils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(serviceName string, zapConfig zap.Config) (*zap.Logger, error) {
	zapConfig.EncoderConfig.TimeKey = "time"
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(timeutils.DefaultTimeFormat))
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.WithOptions(zap.Fields(zap.String("service_name", serviceName)))
	return logger, nil
}

func NewProductionLogger(serviceName string) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	return newLogger(serviceName, zapConfig)
}

func NewDevelopmentLogger(serviceName string) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zapConfig.Sampling = nil
	return newLogger(serviceName, zapConfig)
}

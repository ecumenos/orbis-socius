package httputils

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/ecumenos/orbis-socius/pkg/toolkit/contextutils"
	"github.com/ecumenos/orbis-socius/pkg/toolkit/primitives"
	"github.com/ecumenos/orbis-socius/pkg/toolkit/timeutils"
	"github.com/ecumenos/orbis-socius/schemas/common"
	"go.uber.org/zap"
)

type writer struct {
	rw http.ResponseWriter
	l  *zap.Logger
}

func NewWriter(l *zap.Logger, rw http.ResponseWriter) *writer {
	return &writer{
		rw: rw,
		l:  l,
	}
}

func CreateSuccess[T interface{}](ctx context.Context, payload T) (*common.SuccessResp[T], error) {
	metadata, err := GetMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &common.SuccessResp[T]{
		Data:     payload,
		Metadata: metadata,
		Status:   common.SuccessStatus,
	}, nil
}

func (w *writer) WriteSuccess(ctx context.Context, payload interface{}) error {
	resp, err := CreateSuccess(ctx, payload)
	if err != nil {
		return err
	}
	w.writeHeaders(nil, http.StatusOK)
	return w.write(resp)
}

func CreateFail[T interface{}](ctx context.Context, err error, data T) (*common.FailureResp[T], error) {
	metadata, err := GetMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &common.FailureResp[T]{
		Data:     data,
		Message:  err.Error(),
		Metadata: metadata,
		Status:   common.FailureStatus,
	}, nil
}

func (w *writer) WriteFail(ctx context.Context, inErr error, data interface{}) error {
	if inErr == nil {
		return nil
	}

	resp, err := CreateFail(ctx, inErr, data)
	if err != nil {
		return err
	}
	w.writeHeaders(nil, http.StatusBadRequest)
	return w.write(resp)
}

func CreateError(ctx context.Context, err error) (*common.ErrorResp, error) {
	metadata, err := GetMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &common.ErrorResp{
		Message:  err.Error(),
		Metadata: metadata,
		Status:   common.ErrorStatus,
	}, nil
}

func (w *writer) WriteError(ctx context.Context, inErr error) error {
	if inErr == nil {
		return nil
	}

	resp, err := CreateError(ctx, inErr)
	if err != nil {
		return err
	}
	w.writeHeaders(nil, http.StatusInternalServerError)
	return w.write(resp)
}

func (w *writer) write(payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err := w.rw.Write(b); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeHeaders(headers map[string]string, statusCode int) {
	w.rw.Header().Set("Content-Type", "application/json")
	for key, value := range headers {
		w.rw.Header().Set(key, value)
	}
	w.rw.WriteHeader(statusCode)
}

func WriteSuccess(ctx context.Context, l *zap.Logger, rw http.ResponseWriter, payload interface{}) {
	if err := NewWriter(l, rw).WriteSuccess(ctx, payload); err != nil {
		l.Error("can not write success response", zap.Error(err))
	}
}

func WriteFail(ctx context.Context, l *zap.Logger, rw http.ResponseWriter, inErr error, data interface{}) {
	if err := NewWriter(l, rw).WriteFail(ctx, inErr, data); err != nil {
		l.Error("can not write fail response", zap.Error(err))
	}
}

func WriteError(ctx context.Context, l *zap.Logger, rw http.ResponseWriter, inErr error) {
	if err := NewWriter(l, rw).WriteError(ctx, inErr); err != nil {
		l.Error("can not write error response", zap.Error(err))
	}
}

func GetMetadata(ctx context.Context) (*common.Metadata, error) {
	duration, err := GetRequestDuration(ctx)
	if err != nil {
		return nil, err
	}

	return &common.Metadata{
		RequestID: contextutils.GetValueFromContext(ctx, contextutils.RequestIDKey),
		Timestamp: timeutils.TimeToString(time.Now()),
		Duration:  duration,
		Version:   versioninfo.Short(),
	}, nil
}

func GetRequestDuration(ctx context.Context) (int, error) {
	str := contextutils.GetValueFromContext(ctx, contextutils.StartRequestTimestampKey)
	if str == "" {
		return 0, nil
	}
	start, err := primitives.StringToInt64(str)
	if err != nil {
		return 0, err
	}
	diff := time.Now().UnixNano() - start

	return int(diff), nil
}

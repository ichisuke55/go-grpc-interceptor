package zap

import (
	"bytes"
	"context"
	"testing"

	"github.com/higebu/go-grpc-interceptor/panichandler"
	"github.com/higebu/go-grpc-interceptor/zap/zapctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func NewTestLogger(buf *bytes.Buffer, options ...zap.Option) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(buf), zap.DebugLevel)
	return zap.New(core).WithOptions(options...)
}

func TestUnaryServer(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewTestLogger(buf)

	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test error")
	}

	panichandler.InstallPanicHandler(LogPanicWithStackTrace)

	ctx := context.Background()
	ctx = zapctx.NewContext(ctx, logger)
	_, err := panichandler.UnaryServerInterceptor(ctx, "xyz", unaryInfo, unaryHandler)
	if err == nil {
		t.Fatalf("unexpected success")
	}

	if got, want := grpc.Code(err), codes.Internal; got != want {
		t.Errorf("expect grpc.Code to %s, but got %s", want, got)
	}
}

package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type key string

const (
	loggerKey key = "logger"
)

func CtxWithLog(ctx context.Context, log zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) zerolog.Logger {
	if log := ctx.Value(loggerKey); log != nil {
		return log.(zerolog.Logger)
	}
	return zerolog.Nop()
}

package logger

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	Level string `conf:"default:error"`
	Human bool
}

func New(cfg Config) (zerolog.Logger, error) {
	var output io.Writer = os.Stdout
	if cfg.Human {
		output = zerolog.ConsoleWriter{Out: os.Stderr}
	}
	log := zerolog.New(output)

	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return log, fmt.Errorf("failed to parse log: %w", err)
	}
	log = log.Level(lvl)

	return log, nil
}

func DefaultLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.ErrorLevel)
}

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

func InjectLoggerMiddleware(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			l := log.With().
				Str("url", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Logger()

			r = r.WithContext(CtxWithLog(ctx, l))

			next.ServeHTTP(w, r)
		})
	}
}

package logger

import (
	"net/http"

	"github.com/rs/zerolog"
)

func LoggerMiddleware(log zerolog.Logger) func(http.Handler) http.Handler {
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

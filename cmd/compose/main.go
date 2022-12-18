package main

import (
	"context"
	"time"

	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikenai/gowork/cmd/compose/config"
	"github.com/mikenai/gowork/cmd/compose/handlers"
	"github.com/mikenai/gowork/cmd/compose/pkg/stub"
	"github.com/mikenai/gowork/cmd/compose/pkg/usersapi"
	"github.com/mikenai/gowork/pkg/logger"
)

func main() {
	defaultLog := logger.DefaultLogger()

	cfg, help, err := config.New()
	if err != nil {
		if help != "" {
			defaultLog.Fatal().Msg(help.String())
		}
		defaultLog.Fatal().Err(err).Msg("failed to parse config")
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		defaultLog.Fatal().Err(err).Msg("failed to init logger")
	}

	log.Info().Msg("starting")
	defer log.Info().Msg("shudown")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logger.LoggerMiddleware(log))

	cl := http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     10,
			MaxIdleConns:        5,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     time.Second * 1,
		},
		Timeout: time.Second * 2,
	}

	stub := &stub.Client{
		BaseURL: "http://localhost:9090",
		Http:    cl,
	}

	users := &usersapi.Client{
		BaseURL: "localhost:8080",
		Http:    cl,
	}

	h := handlers.Handler{
		PostsAPI:    stub,
		ProfilesAPI: stub,
		UsersAPI:    users,

		Log: log,
	}

	r.Get("/{user_id}", h.UserPage)

	s := http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: r,

		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,

		IdleTimeout: cfg.HTTP.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Info().Msg("signal received")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefullTimeout)
		defer cancel()

		log.Info().Msg("shutting down")
		s.Shutdown(ctx)
	}()

	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to close the server")
	}
}

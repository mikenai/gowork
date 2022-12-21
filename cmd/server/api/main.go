package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mikenai/gowork/cmd/server/config"
	"github.com/mikenai/gowork/internal/handlers"
	userstorage "github.com/mikenai/gowork/internal/storage/users"
	"github.com/mikenai/gowork/internal/users"
	"github.com/mikenai/gowork/pkg/dbcollector"
	"github.com/mikenai/gowork/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	db, err := sql.Open("sqlite3", cfg.DB.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close() // always close resources

	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

	ur := userstorage.New(db)
	us := users.New(ur)
	uh := handlers.NewUsers(us)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logger.LoggerMiddleware(log))

	prometheus.MustRegister(dbcollector.NewSQLDatabaseCollector("general", "main", "sqlite", db))
	r.Mount("/metrics", promhttp.Handler())

	r.Mount("/users", uh.Routes())

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
		fmt.Println(s.ListenAndServe())
		stop()
	}()

	<-ctx.Done()
	log.Info().Msg("signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefullTimeout)
	defer cancel()

	log.Info().Msg("shutting down")
	s.Shutdown(ctx)
}

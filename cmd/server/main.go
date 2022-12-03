package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("starting")
	defer fmt.Println("shutdown")

	cfg, help, err := config.New()
	if err != nil {
		if help != "" {
			log.Fatal(help)
		}
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", cfg.DB.DSN)
	if err != nil {
		log.Fatal("failed to connect to database", err)
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

	go func() {
		fmt.Println(s.ListenAndServe())
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	fmt.Println("signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefullTimeout)
	defer cancel()

	fmt.Println("shutting down")
	s.Shutdown(ctx)
}

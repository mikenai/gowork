package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
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

	db, err := sql.Open("sqlite3", "./tmp/db.sqlite3")
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(10)

	db.SetConnMaxLifetime(time.Second * 5)
	db.SetConnMaxIdleTime(time.Second * 1)

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
		Addr:    ":8080",
		Handler: r,

		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,

		IdleTimeout: 2 * time.Second,
	}

	go func() {
		fmt.Println(s.ListenAndServe())
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	fmt.Println("signal received")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Println("shutting down")
	s.Shutdown(ctx)
}

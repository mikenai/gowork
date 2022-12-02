package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikenai/gowork/internal/handlers"
	"github.com/mikenai/gowork/internal/users"
)

func main() {
	fmt.Println("starting")
	defer fmt.Println("shutdown")

	us := users.Service{}

	uh := handlers.NewUsers(us)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/users", uh.Create)
	r.Get("/users/{id}", uh.Fetch)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		fmt.Println("over")
	})

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

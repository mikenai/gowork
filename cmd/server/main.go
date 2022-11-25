package main

import (
	"fmt"
	"net/http"

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

	s := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println(s.ListenAndServe())
}

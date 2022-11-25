package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikenai/gowork/internal/handlers"
)

func main() {
	uh := handlers.Users{}

	r := chi.NewRouter()

	r.Post("/users", uh.Create)

	s := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println(s.ListenAndServe())
}

package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	stub "github.com/mikenai/gowork/cmd/stub/handlers"
	"github.com/rs/zerolog"
)

type config struct {
	HTTPAddr string `conf:"default::9090"`
}

func main() {
	log := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	log.Info().Msg("starting")
	defer log.Info().Msg("shutdown")

	cfg := config{}

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			log.Fatal().Msg(help)
		}
		log.Fatal().Err(err).Msg("failed to parse log")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/profiles/{user_id}", stub.ProfileHandler)
	r.Get("/{user_id}/posts", stub.PostsHandler)

	log.Info().Str("http_addr", cfg.HTTPAddr).Msg("server starting")
	http.ListenAndServe(cfg.HTTPAddr, r)
}

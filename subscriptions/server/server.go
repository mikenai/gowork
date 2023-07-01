package server

import (
	"context"
	"fmt"
	"main/client/database"
	"main/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
)

// Server holds the HTTP server, router, config and all clients.
type Server struct {
	Config *config.Config
	HTTP   *http.Server
	Router chi.Router
	DB     *database.Client
}

// Create sets up the HTTP server, router and all clients.
// Returns an error if an error occurs.
func (s *Server) Create(ctx context.Context, config *config.Config) error {
	var dbClient database.Client

	if err := dbClient.Init(ctx, config); err != nil {
		return fmt.Errorf("initialize db client: %w", err)
	}

	s.DB = &dbClient
	s.Config = config
	s.Router = chi.NewRouter()
	s.HTTP = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Config.Port),
		Handler: s.Router,
	}

	s.setupRoutes()

	return nil
}

// Serve tells the server to start listening and serve HTTP requests.
// It also makes sure that the server gracefully shuts down on exit.
// Returns an error if an error occurs.
func (s *Server) Serve(ctx context.Context) error {
	go func(ctx context.Context, s *Server) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		log.Info("Shutdown signal received")

		if err := s.HTTP.Shutdown(ctx); err != nil {
			log.Error(err.Error())
		}
		if err := s.DB.Close(); err != nil {
			log.Error(err.Error())
		}
	}(ctx, s)

	log.Infof("Ready at: %s", s.Config.Port)

	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}
	return nil
}

package main

import (
	"context"
	"main/config"
	"main/server"

	"github.com/dbmedialab/pkg/logrustic"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(logrustic.NewFormatter("Subscriptions."))

	log.Info("Starting ...")

	ctx := context.Background()
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	var s server.Server
	if err := s.Create(ctx, &config); err != nil {
		log.Fatal(err.Error())
	}

	if err := s.Serve(ctx); err != nil {
		log.Fatal(err.Error())
	}

}

package main

import (
	"database/sql"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikenai/gowork/cmd/server/config"
	handlers "github.com/mikenai/gowork/internal/handlers/grpc"
	pb "github.com/mikenai/gowork/internal/proto"
	userstorage "github.com/mikenai/gowork/internal/storage/users"
	"github.com/mikenai/gowork/internal/users"
	"github.com/mikenai/gowork/pkg/logger"
	"google.golang.org/grpc"
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
	defer log.Info().Msg("shutdown")

	db, err := sql.Open("sqlite3", cfg.DB.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close() // always close resources

	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

	storage := userstorage.New(db)
	users := users.New(storage)

	lis, err := net.Listen("tcp", cfg.GRCPPort)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed listen on port %s", cfg.GRCPPort)
	}
	defer lis.Close()

	s := grpc.NewServer()
	// TODO: pass logger in context
	pb.RegisterUsersServiceServer(s, handlers.NewUsersGRPC(users, log))
	defer s.Stop()

	log.Info().Msg("grpc server started")
	defer log.Info().Msg("shutdown")

	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}

	// TODO: add gracefully shutdown
}

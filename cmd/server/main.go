package main

import (
	"database/sql"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikenai/gowork/cmd/server/config"
	"github.com/mikenai/gowork/internal/handlers"
	"github.com/mikenai/gowork/internal/shared/protobuf"
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

	port := cfg.HTTP.Addr

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal().Msg("failed to listen")
	}

	gsrv := grpc.NewServer()
	protobuf.RegisterUsersServer(gsrv, &uh)
	if err := gsrv.Serve(lis); err != nil {
		log.Fatal().Msg("failed to serve")
	}
}

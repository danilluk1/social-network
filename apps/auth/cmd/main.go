package main

import (
	"context"
	"fmt"
	"net"
	"os"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/grpc_impl"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/danilluk1/social-network/libs/grpc/servers"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err)
	}

	if cfg.AppEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := pgx.Connect(context.Background(), cfg.AuthPostgresUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database:")
	}
	defer conn.Close(context.Background())

	runDBMigation(cfg.AuthPostgresUrl, config.DBSource)

	store := db.NewStore(conn)

	runGrpcServer(cfg, store)
}

func runDBMigation(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migration instance:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up:")
	}

	log.Info().Msg("db migrate successfully")
}

func runGrpcServer(config *config.Config, store db.Store) {
	server, err := grpc_impl.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(grpc_impl.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	auth.RegisterAuthServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.AUTH_SERVER_PORT))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("🚀 start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server:")
	}

}

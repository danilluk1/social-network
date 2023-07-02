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
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

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

	store := db.NewStore(conn)

	grpcServerImpl, err := grpc_impl.NewServer(cfg, store)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create auth grpc server:")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.AUTH_SERVER_PORT))
	if err != nil {
		log.Error().Err(err).Msg("Failed to listen:")
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, grpcServerImpl)
	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	//Write support of logrus instead of zap inside fx
	app := fx.New(
		fx.Provide(
			func(lc fx.Lifecycle) *pgx.Conn {
				lc.Append(fx.Hook{
					OnStop: func(context.Context) error {
						return conn.Close(ctx)
					},
				})
				return conn
			},
			func() *config.Config { return cfg },
		),
	)

	log.Info().Msg("Auth microservice started")
	app.Run()
}

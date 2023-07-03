package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/grpc_impl"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/danilluk1/social-network/libs/grpc/servers"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	zap.ReplaceGlobals(logger)

	conn, err := pgx.Connect(ctx, cfg.AuthPostgresUrl)
	if err != nil {
		logger.Sugar().Error(err)
	}
	defer conn.Close(ctx)

	store := db.NewStore(conn)

	grpcServerImpl, err := grpc_impl.NewServer(cfg, store, logger)
	if err != nil {
		logger.Sugar().Error("Failed to create auth grpc server:", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.AUTH_SERVER_PORT))
	if err != nil {
		logger.Sugar().Error("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, grpcServerImpl)
	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	logger.Sugar().Info("Auth microservice started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Sugar().Info("Exiting...")
	cancel()
}

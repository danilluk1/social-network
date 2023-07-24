package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/danilluk1/social-network/apps/auth/internal/conf"
	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/gapi"
	"github.com/danilluk1/social-network/apps/auth/internal/token"
	"github.com/danilluk1/social-network/apps/auth/internal/utilities"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/danilluk1/social-network/libs/grpc/servers"
	"github.com/danilluk1/social-network/libs/kafka/topics"
	"github.com/jackc/pgx/v5"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var serveCmd = cobra.Command{
	Long: "Start GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func serve(ctx context.Context) {
	config, err := conf.LoadGlobal(configFile)
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger
	if config.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	conn, err := pgx.Connect(ctx, config.DB.URL)
	if err != nil {
		logger.Sugar().Fatalf("error openning database: %+v", err)
	}
	defer conn.Close(ctx)

	store := db.NewStore(conn)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Kafka.KafkaUrl),
		Topic:                  topics.Mail,
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()
	schemaClient := srclient.CreateSchemaRegistryClient(config.Kafka.SchemaRegistryUrl)
	schema, err := schemaClient.GetLatestSchema(topics.Mail)
	if schema == nil || err == nil {
		schemaBytes, err := ioutil.ReadFile(config.Kafka.SchemasPath + topics.Mail + ".avsc")
		if err != nil {
			logger.Sugar().Fatalf("failed to read schemas path: %+v", err)
		}
		_, err = schemaClient.CreateSchema(topics.Mail, string(schemaBytes), srclient.Avro)
		if err != nil {
			logger.Sugar().Fatalf("failed to create schemas: %+v", err)
		}
	}

	tokenMaker, err := token.NewPasetoMaker(config.PASETO.Secret)
	if err != nil {
		logger.Sugar().Fatalf("Failed to create token maker: %+v", err)
	}

	services := &gapi.Services{
		Conf:         config,
		Store:        store,
		TokenMaker:   tokenMaker,
		KafkaWriter:  writer,
		SchemaClient: schemaClient,
	}

	grpcApi := gapi.NewGAPIWithVersion(ctx, services, utilities.Version)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.AUTH_SERVER_PORT))
	if err != nil {
		logger.Sugar().Fatalf("failed to listen: %+v", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLoggerWrapper(logger))
	grpcServer := grpc.NewServer(grpcLogger)
	auth.RegisterAuthServer(grpcServer, grpcApi)
	reflection.Register(grpcServer)

	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	addr := net.JoinHostPort(config.GRPC.Host, config.GRPC.Port)
	logger.Sugar().Infof("Auth microservice started on: %s", addr)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Sugar().Info("Exiting...")
}

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/danilluk1/social-network/apps/auth/internal/conf"
	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/gapi"
	"github.com/danilluk1/social-network/apps/auth/internal/token"
	"github.com/danilluk1/social-network/apps/auth/internal/utilities"
	"github.com/danilluk1/social-network/libs/grpc/servers"
	"github.com/danilluk1/social-network/libs/kafka/topics"
	"github.com/jackc/pgx/v5"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func serve(ctx context.Context) {
	config, err := conf.LoadGlobal(configFile)
	if err != nil {
		logrus.WithError(err).Fatal("unable to load config")
	}

	conn, err := pgx.Connect(ctx, config.DB.URL)
	if err != nil {
		logrus.Fatalf("error opening database: %+v", err)
	}
	defer conn.Close(ctx)

	store := db.NewStore(conn)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.KafkaUrl),
		Topic:                  topics.Mail,
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()
	schemaClient := srclient.CreateSchemaRegistryClient(config.SchemaRegistryUrl)
	schema, err := schemaClient.GetLatestSchema(topics.Mail)
	if schema == nil || err == nil {
		schemaBytes, err := ioutil.ReadFile(cfg.SchemasPath + topics.Mail + ".avsc")
		if err != nil {
			logrus.Fatalf("failed to read schemas path: %+v", err)
		}
		_, err = schemaClient.CreateSchema(topics.Mail, string(schemaBytes), srclient.Avro)
		if err != nil {
			logrus.Fatalf("failed to create schemas: %+v", err)
		}
	}

	tokenMaker, err := token.NewPasetoMaker(config.PASETO.Secret)
	if err != nil {
		logrus.Fatalf("failed to create token maker: %+v", err)
	}

	services := &gapi.Services{
		Conf:         config,
		Store:        store,
		TokenMaker:   tokenMaker,
		KafkaWriter:  writer,
		SchemaClient: schemaClient,
	}

	grpcApi := gapi.NewWithVersion(ctx, services, utilities.Version)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.AUTH_SERVER_PORT))
	if err != nil {
		logrus.Fatalf("Failed to listen: %+v", err)
	}

	addr := net.JoinHostPort(config.GRPC.Host, config.GRPC.Port)
	logrus.Infof("Auth API started on: %s", addr)

	grpcApi.ListenAndServe(ctx, addr)
}

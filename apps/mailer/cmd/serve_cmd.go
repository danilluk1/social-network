package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danilluk1/social-network/apps/mailer/internal/conf"
	"github.com/danilluk1/social-network/apps/mailer/internal/mail"
	"github.com/danilluk1/social-network/apps/mailer/internal/observability"
	"github.com/danilluk1/social-network/libs/kafka/topics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var serveCmd = cobra.Command{
	Long: "Start mailer service",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func serve(ctx context.Context) {
	cfg, err := conf.Load(configFile)
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

	if cfg.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	zap.ReplaceGlobals(logger)

	shutdown, err := observability.InitProviderWithJaegerExporter(ctx, cfg)
	if err != nil {
		logger.Sugar().Fatalf("Failed to initialize opentelemtry provider: %v", err)
	}
	defer shutdown(ctx)

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Kafka.KafkaUrl},
		Topic:   topics.Mail,
		GroupID: cfg.Kafka.GroupID,
	})
	defer kafkaReader.Close()

	services := &mail.Services{
		Logger:         logger,
		Mail:           gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass),
		Reader:         kafkaReader,
		SchemaRegistry: srclient.CreateSchemaRegistryClient(cfg.Kafka.SchemaRegistryUrl),
	}
	reader := mail.New(services)
	go reader.Start(ctx, cfg)

	logger.Sugar().Info("Mailer microservice started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Sugar().Info("Exiting...")
}

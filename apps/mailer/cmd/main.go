package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danilluk1/social-network/apps/mailer/internal/mail"
	"github.com/danilluk1/social-network/apps/mailer/internal/types"
	"github.com/danilluk1/social-network/libs/config"
	topics "github.com/danilluk1/social-network/libs/kafka/topics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if cfg.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
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

	services := &types.Services{
		Logger: logger,
		Mail:   gomail.NewDialer(cfg.MailHost, cfg.MailPort, cfg.MailUser, cfg.MailPass),
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.KafkaUrl},
			Topic:   topics.Mail,
			GroupID: "mailer",
		}),
		SchemaRegistry: srclient.CreateSchemaRegistryClient(cfg.SchemaRegistryUrl),
	}

	reader := mail.New(services)
	go reader.Start(ctx)

	logger.Sugar().Info("Mailer microservice started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Sugar().Info("Exiting...")
	cancel()
}

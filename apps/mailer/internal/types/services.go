package types

import (
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type Services struct {
	Logger         *zap.Logger
	Mail           *gomail.Dialer
	Reader         *kafka.Reader
	SchemaRegistry *srclient.SchemaRegistryClient
}

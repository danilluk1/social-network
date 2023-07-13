package grpc_impl

import (
	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/token"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Server struct {
	auth.UnimplementedAuthServer
	config         *config.Config
	store          db.Store
	logger         *zap.Logger
	tokenMaker     token.Maker
	kafkaWriter    *kafka.Writer
	schemaRegistry *srclient.SchemaRegistryClient
}

func NewServer(config *config.Config, store db.Store, logger *zap.Logger, kafka *kafka.Writer, registry *srclient.SchemaRegistryClient) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.JwtSecret)
	if err != nil {
		return nil, err
	}

	return &Server{
		config:         config,
		store:          store,
		tokenMaker:     tokenMaker,
		logger:         logger,
		kafkaWriter:    kafka,
		schemaRegistry: registry,
	}, nil
}

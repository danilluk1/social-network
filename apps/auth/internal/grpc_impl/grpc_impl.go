package grpc_impl

import (
	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/token"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"go.uber.org/zap"
)

type Server struct {
	auth.UnimplementedAuthServer
	config     *config.Config
	store      db.Store
	logger     *zap.Logger
	tokenMaker token.Maker
}

func NewServer(config *config.Config, store db.Store, logger *zap.Logger) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.JwtSecret)
	if err != nil {
		return nil, err
	}

	return &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		logger:     logger,
	}, nil
}

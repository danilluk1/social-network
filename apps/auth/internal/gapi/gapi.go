package gapi

import (
	"context"

	"github.com/danilluk1/social-network/apps/auth/internal/conf"
	db "github.com/danilluk1/social-network/apps/auth/internal/db/sqlc"
	"github.com/danilluk1/social-network/apps/auth/internal/token"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

const (
	defaultVersion = "unknown version"
)

type GAPI struct {
	auth.UnimplementedAuthServer
	services *Services
	version  string
}

type Services struct {
	Conf         *conf.GlobalConfiguration
	Store        db.Store
	TokenMaker   token.Maker
	KafkaWriter  *kafka.Writer
	SchemaClient *srclient.SchemaRegistryClient
}

func NewGAPI(services *Services) *GAPI {
	return NewGAPIWithVersion(context.Background(), services, defaultVersion)
}

// func (g *GAPI) deprecationNotices(ctx context.Context) {
// 	config := g.services.Conf

// 	log := logrus.WithField("component", "gapi")

// }

func NewGAPIWithVersion(ctx context.Context, services *Services, version string) *GAPI {
	gapi := &GAPI{
		services: services,
		version:  version,
	}

	return gapi
}

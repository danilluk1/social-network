package gapi

import (
	"context"

	"github.com/danilluk1/social-network/apps/auth/internal/shared/contracts"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

const (
	defaultVersion = "unknown version"
)

var grpcMetricsAttr = api.WithAttributes(
	attribute.Key("MetricsType").String("Http"),
)

type GAPI struct {
	authMetrics *contracts.AuthMetrics
	auth.UnimplementedAuthServer
}

func NewGAPI() *GAPI {
	return NewGAPIWithVersion(context.Background(), defaultVersion)
}

func NewGAPIWithVersion(ctx context.Context, version string) *GAPI {
	gapi := &GAPI{}

	return gapi
}

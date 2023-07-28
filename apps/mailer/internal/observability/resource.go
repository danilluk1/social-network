package observability

import (
	"context"
	"log"
	"os"

	"github.com/danilluk1/social-network/apps/mailer/internal/conf"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func newResource(ctx context.Context, conf *conf.Configuration) *resource.Resource {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("mailer"),
			attribute.String("environment",
				os.Getenv("GO_ENV")),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	return res
}

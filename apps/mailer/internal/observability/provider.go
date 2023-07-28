package observability

import (
	"context"
	"log"

	"github.com/danilluk1/social-network/apps/mailer/internal/conf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitProviderWithJaegerExporter(ctx context.Context, conf *conf.Configuration) (func(context.Context) error, error) {
	exp, err := exporterToJagger(conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(getSampler(conf)),
		trace.WithBatcher(exp),
		trace.WithResource(newResource(ctx, conf)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

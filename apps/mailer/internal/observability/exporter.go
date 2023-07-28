package observability

import (
	"github.com/danilluk1/social-network/apps/mailer/internal/conf"
	"go.opentelemetry.io/otel/exporters/jaeger"
)

func exporterToJagger(conf *conf.Configuration) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.OpenTelemetryCollectorUrl)))
}

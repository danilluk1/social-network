package observability

import (
	"github.com/danilluk1/social-network/apps/mailer/internal/conf"
	"go.opentelemetry.io/otel/sdk/trace"
)

func getSampler(conf *conf.Configuration) trace.Sampler {
	switch conf.AppEnv {
	case "development":
		return trace.AlwaysSample()
	case "production":
		return trace.ParentBased(trace.TraceIDRatioBased(0.5))
	default:
		return trace.AlwaysSample()
	}
}

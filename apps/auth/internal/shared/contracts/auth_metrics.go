package contracts

import "go.opentelemetry.io/otel/metric"

type AuthMetrics struct {
	SignUpRequests metric.Float64Counter
	LoginRequests  metric.Float64Counter
}

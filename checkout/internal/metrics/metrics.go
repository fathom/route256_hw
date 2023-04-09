package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "checkout",
		Subsystem: "grpc",
		Name:      "requests_total",
	})

	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "checkout",
		Subsystem: "grpc",
		Name:      "responses_total",
	},
		[]string{"status", "handler"},
	)

	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "checkout",
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status", "handler"},
	)

	HistogramClientTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "checkout",
		Subsystem: "grpc",
		Name:      "histogram_client_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status", "handler"},
	)
)

func NewMetricHandler() http.Handler {
	return promhttp.Handler()
}

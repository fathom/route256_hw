package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	JobCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "notifications",
		Subsystem: "kafka",
		Name:      "job_total",
	})
)

func NewMetricHandler() http.Handler {
	return promhttp.Handler()
}

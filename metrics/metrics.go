package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns the global http.Handler that provides the prometheus
// metrics format on GET requests.
func Handler() http.Handler {
	return promhttp.Handler()
}

// NewCounter initiates and registers a metrics counter.
func NewCounterVec(name string, desc string, labels []string) *prometheus.CounterVec {
	c := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: desc,
	}, labels)
	return c
}

// NewHistogramVec initiates and registers a latency histogram.
func NewHistogramVec(name string, desc string, buckets []float64, labels ...string) *prometheus.HistogramVec {
	hv := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    desc,
			Buckets: buckets,
		},
		labels,
	)
	return hv
}

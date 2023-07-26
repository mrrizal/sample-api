package observer

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ActiveConnections  *prometheus.Gauge
	RequestsTotal      *prometheus.CounterVec
	DurationHistorgram *prometheus.HistogramVec
	DurationSummary    *prometheus.SummaryVec
}

func InitMetrics(reg prometheus.Registerer) *Metrics {
	activeConnections := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of acticve connections",
		},
	)

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests received.",
		},
		[]string{"method", "endpoint"},
	)

	durationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	durationSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_summary_seconds",
			Help:       "Summary of HTTP request durations.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "endpoint", "status_code"},
	)

	metrics := &Metrics{
		ActiveConnections:  &activeConnections,
		RequestsTotal:      requestsTotal,
		DurationHistorgram: durationHistogram,
		DurationSummary:    durationSummary,
	}

	reg.MustRegister(*metrics.ActiveConnections)
	reg.MustRegister(metrics.RequestsTotal)
	reg.MustRegister(metrics.DurationHistorgram)
	reg.MustRegister(metrics.DurationSummary)
	return metrics
}

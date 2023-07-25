package observer

import "github.com/prometheus/client_golang/prometheus"

var (
	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections.",
		},
	)

	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests received.",
		},
		[]string{"method", "endpoint"},
	)

	HTTPDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	HTTPDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_summary_seconds",
			Help:       "Summary of HTTP request durations.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "status_code"},
	)
)

func InitMetric() {
	prometheus.MustRegister(ActiveConnections)
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPDurationHistogram)
	prometheus.MustRegister(HTTPDurationSummary)
}

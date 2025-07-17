package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Define Prometheus metrics
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
)

// Use these metrics in your handlers:
// httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
// httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())

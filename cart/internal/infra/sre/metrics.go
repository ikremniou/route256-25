package sre

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalHttpRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_http_total_requests",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "code"},
	)
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cart_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "code"},
	)
	TotalExternalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_external_total_requests",
			Help: "Total number of external requests",
		},
		[]string{"action", "status"},
	)
	ExternalRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cart_external_request_duration_seconds",
			Help:    "Duration of external requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"action", "status"},
	)
	InMemoryCartItems = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cart_in_memory_items_count",
			Help: "Number of carts in memory",
		},
		[]string{},
	)
)

func TrackExternalRequest(action string, err error, startTime time.Time) {
	duration := time.Since(startTime)

	status := "success"
	if err != nil {
		status = "error"
	}

	TotalExternalRequests.With(prometheus.Labels{
		"action": action,
		"status": status,
	}).Inc()
	ExternalRequestDuration.With(prometheus.Labels{
		"action": action,
		"status": status,
	}).Observe(duration.Seconds())
}

func TrackHttpRequest(method, path string, statusCode int, startTime time.Time) {
	duration := time.Since(startTime)

	TotalHttpRequests.With(prometheus.Labels{
		"method": method,
		"path":   path,
		"code":   strconv.Itoa(statusCode),
	}).Inc()
	HttpRequestDuration.With(prometheus.Labels{
		"method": method,
		"path":   path,
		"code":   strconv.Itoa(statusCode),
	}).Observe(duration.Seconds())
}

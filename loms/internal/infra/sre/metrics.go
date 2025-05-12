package sre

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/status"
)

var (
	TotalGrpcRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loms_grpc_total_requests",
			Help: "Total number of requests",
		},
		[]string{"method", "code"},
	)
	GrpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loms_grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)
	TotalHttpRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loms_http_total_requests",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "code"},
	)
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loms_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "code"},
	)
	TotalDatabaseRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loms_database_total_requests",
			Help: "Total number of database requests",
		},
		[]string{"action", "category", "status"},
	)
	DatabaseRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loms_database_request_duration_seconds",
			Help:    "Duration of database requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"action", "category", "status"},
	)
	TotalExternalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loms_external_total_requests",
			Help: "Total number of external requests",
		},
		[]string{"action", "status"},
	)
	ExternalRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loms_external_request_duration_seconds",
			Help:    "Duration of external requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"action", "status"},
	)
)

func TrackDbRequest(action, category string, err error, start time.Time) {
	duration := time.Since(start)
	status := "success"
	if err != nil {
		status = "error"
	}

	TotalDatabaseRequests.With(prometheus.Labels{
		"action":   action,
		"category": category,
		"status":   status,
	}).Inc()

	DatabaseRequestDuration.With(prometheus.Labels{
		"action":   action,
		"category": category,
		"status":   status,
	}).Observe(duration.Seconds())
}

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

func TrackGrpcRequest(method string, err error, startTime time.Time) {
	duration := time.Since(startTime)
	statusCode := status.Code(err).String()

	TotalGrpcRequests.With(prometheus.Labels{
		"method": method,
		"code":   statusCode,
	}).Inc()
	GrpcRequestDuration.With(prometheus.Labels{
		"method": method,
		"code":   statusCode,
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

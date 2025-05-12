package middleware

import (
	"net/http"
	"route256/cart/internal/infra/logger"
	"time"
)

type RequestLoggerMiddleware struct {
	mux http.Handler
}

func NewRequestLoggerMiddleware(mux http.Handler) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{mux: mux}
}

func (request *RequestLoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writerWithLogger := newWriterLogger(w)

	start := time.Now()
	request.mux.ServeHTTP(writerWithLogger, r)
	duration := time.Since(start)

	logger.Debug("Processed request", "method", r.Method, "uri", r.RequestURI, "proto", r.Proto,
		"duration", duration.String(), "status", writerWithLogger.statusCode, "bytes", writerWithLogger.bytes)
}

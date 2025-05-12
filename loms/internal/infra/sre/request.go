package sre

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (ww *ResponseWriterWrapper) WriteHeader(code int) {
	ww.statusCode = code
	ww.ResponseWriter.WriteHeader(code)
}

func GrpcMw(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx, span := otel.Tracer("request").Start(ctx, info.FullMethod)
	defer span.End()

	startTime := time.Now()
	result, err := handler(ctx, req)
	TrackGrpcRequest(info.FullMethod, err, startTime)

	return result, err
}

func NewHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var spanName = fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		ctx, span := otel.Tracer("request").Start(r.Context(), spanName)
		defer span.End()

		ww := &ResponseWriterWrapper{ResponseWriter: w}
		r = r.WithContext(ctx)

		startTime := time.Now()
		h.ServeHTTP(ww, r)
		TrackHttpRequest(r.Method, r.URL.Path, ww.statusCode, startTime)
	})
}

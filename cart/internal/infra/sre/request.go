package sre

import (
	"net/http"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (ww *ResponseWriterWrapper) WriteHeader(code int) {
	ww.statusCode = code
	ww.ResponseWriter.WriteHeader(code)
}

func NewHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &ResponseWriterWrapper{ResponseWriter: w}

		startTime := time.Now()
		h.ServeHTTP(ww, r)
		TrackHttpRequest(r.Method, r.URL.Path, ww.statusCode, startTime)
	})
}

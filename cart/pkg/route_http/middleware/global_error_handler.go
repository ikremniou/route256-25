package middleware

import (
	"fmt"
	"net/http"
	"route256/cart/internal/infra/logger"
	"runtime/debug"
)

type GlobalRequestErrorMiddleware struct {
	mux http.Handler
}

func NewGlobalRequestErrorMiddleware(mux http.Handler) *GlobalRequestErrorMiddleware {
	return &GlobalRequestErrorMiddleware{mux: mux}
}

func (request *GlobalRequestErrorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Panic in request handler", "error", err, "stack", string(debug.Stack()))
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
	}()

	request.mux.ServeHTTP(w, r)
}

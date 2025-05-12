package app

import (
	"net/http"
	"route256/loms/internal/infra/loms_config"
)

func enableCors(h http.Handler, config *loms_config.Config) http.Handler {
	if !config.Server.AllowSwagger {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8086")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

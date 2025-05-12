package route_http

import (
	"encoding/json"
	"net/http"
	"route256/cart/internal/infra/logger"
)

type errorRepose struct {
	Error string `json:"error"`
}

func WriteJson(w http.ResponseWriter, status int, anything any) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(anything)
	if err != nil {
		logger.Error("Unable to encode json response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteErrorJson(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	err = json.NewEncoder(w).Encode(errorRepose{Error: err.Error()})
	if err != nil {
		logger.Error("Unable to encode json error response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

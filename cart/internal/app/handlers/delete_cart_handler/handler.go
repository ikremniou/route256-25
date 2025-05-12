package delete_cart_handler

import (
	"context"
	"net/http"
	"route256/cart/internal/infra/infra_http"
	"route256/cart/pkg/route_http"

	"go.opentelemetry.io/otel"
)

type CartService interface {
	DeleteAll(ctx context.Context, userId int64) error
}

type ClearCartHandler struct {
	service CartService
}

func New(service CartService) *ClearCartHandler {
	return &ClearCartHandler{
		service: service,
	}
}

func (h *ClearCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("handler").Start(r.Context(), "delete_cart_handler.ServeHTTP")
	defer span.End()

	userId, err := infra_http.GetInt64PathValueGt0(r, "user_id")
	if err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	if err := h.service.DeleteAll(ctx, userId); err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

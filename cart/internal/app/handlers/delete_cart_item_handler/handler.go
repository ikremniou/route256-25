package delete_cart_item_handler

import (
	"context"
	"net/http"
	"route256/cart/internal/infra/infra_http"
	"route256/cart/pkg/route_http"

	"go.opentelemetry.io/otel"
)

type CartService interface {
	DeleteBySkuId(ctx context.Context, userId int64, skuId int64) error
}

type DeleteCartItemHandler struct {
	service CartService
}

func New(service CartService) *DeleteCartItemHandler {
	return &DeleteCartItemHandler{
		service: service,
	}
}

func (handler *DeleteCartItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("handler").Start(r.Context(), "delete_cart_item_handler.ServeHTTP")
	defer span.End()

	userId, err := infra_http.GetInt64PathValueGt0(r, "user_id")
	if err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	skuId, err := infra_http.GetInt64PathValueGt0(r, "sku_id")
	if err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	if err := handler.service.DeleteBySkuId(ctx, userId, skuId); err != nil {
		route_http.WriteErrorJson(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

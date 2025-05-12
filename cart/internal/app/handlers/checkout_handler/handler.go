package checkout_handler

import (
	"context"
	"errors"
	"net/http"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/infra_http"
	"route256/cart/pkg/route_http"

	"go.opentelemetry.io/otel"
)

type CartService interface {
	Checkout(ctx context.Context, userId int64) (int64, error)
}

type CheckoutHandler struct {
	cartService CartService
}

func New(cartService CartService) *CheckoutHandler {
	return &CheckoutHandler{cartService: cartService}
}

func (h *CheckoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("handler").Start(r.Context(), "checkout_handler.ServeHTTP")
	defer span.End()

	userId, err := infra_http.GetInt64PathValueGt0(r, "user_id")
	if err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	orderId, err := h.cartService.Checkout(ctx, userId)
	if err != nil {
		if errors.Is(err, model.ErrCreateOrderPreconditionFailed) {
			route_http.WriteErrorJson(w, http.StatusPreconditionFailed, err)

			return
		}

		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	route_http.WriteJson(w, http.StatusOK, &CheckoutResponse{OrderId: orderId})
}

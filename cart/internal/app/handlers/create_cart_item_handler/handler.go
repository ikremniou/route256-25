package create_cart_item_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/infra_http"
	"route256/cart/pkg/route_http"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
)

type CartService interface {
	Create(ctx context.Context, cartItem *model.CartItemModel) (bool, error)
}

type CreateCartItemHandler struct {
	service   CartService
	validator *validator.Validate
}

func New(service CartService) *CreateCartItemHandler {
	return &CreateCartItemHandler{
		service:   service,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (handler *CreateCartItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("handler").Start(r.Context(), "create_cart_item_handler.ServeHTTP")
	defer span.End()
	defer r.Body.Close()

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

	var createItemRequest CreateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&createItemRequest); err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	if err := handler.validator.Struct(&createItemRequest); err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	_, err = handler.service.Create(ctx, &model.CartItemModel{
		UserId: userId,
		SkuId:  skuId,
		Count:  createItemRequest.Count,
	})
	if err != nil {
		if errors.Is(err, model.ErrProductDoesNotExist) || errors.Is(err, model.ErrNotEnoughItemsInStock) {
			route_http.WriteErrorJson(w, http.StatusPreconditionFailed, err)
		} else {
			route_http.WriteErrorJson(w, http.StatusBadRequest, err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

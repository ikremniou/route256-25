package get_cart_items_handler

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
	GetAllOrderBySku(ctx context.Context, userId int64) (model.AllCartItemsModel, error)
}

type GetCartItemBySkuHandler struct {
	service CartService
}

func New(service CartService) *GetCartItemBySkuHandler {
	return &GetCartItemBySkuHandler{service: service}
}

func (handler *GetCartItemBySkuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("handler").Start(r.Context(), "get_cart_items_handler.ServeHTTP")
	defer span.End()

	userId, err := infra_http.GetInt64PathValueGt0(r, "user_id")
	if err != nil {
		route_http.WriteErrorJson(w, http.StatusBadRequest, err)

		return
	}

	items, err := handler.service.GetAllOrderBySku(ctx, userId)
	if err != nil {
		if errors.Is(err, model.ErrTheCartIsEmpty) {
			route_http.WriteErrorJson(w, http.StatusNotFound, err)
		} else {
			route_http.WriteErrorJson(w, http.StatusBadRequest, err)
		}

		return
	}

	var responseItems = make([]GetCartItemResponseItem, len(items.Items))
	for i, item := range items.Items {
		responseItems[i] = GetCartItemResponseItem{
			Sku:   item.SkuId,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		}
	}

	route_http.WriteJson(w, http.StatusOK, GetCartItemsResponse{
		Items:      responseItems,
		TotalPrice: items.Total,
	})
}

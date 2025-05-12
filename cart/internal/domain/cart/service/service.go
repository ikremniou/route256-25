package service

import (
	"context"
	"fmt"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/logger"
	"route256/cart/pkg/route_err_group"
	"sync/atomic"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
)

//go:generate minimock -i ProductService,CartRepository,StocksClient,OrdersClient -p service_test,service_test,service_test,service_test
type CartRepository interface {
	CreateItem(ctx context.Context, item *model.CartItemModel) (bool, error)
	GetAllOrderBySku(ctx context.Context, userId int64) []model.CartItemModel
	DeleteAll(ctx context.Context, userId int64)
	DeleteBySku(ctx context.Context, userId int64, skuId int64)
}

type ProductService interface {
	IsProductExists(ctx context.Context, skuId int64) (bool, error)
	GetProduct(ctx context.Context, skuId int64) (model.ProductModel, error)
	GetProductsAot(ctx context.Context, count int64, startSkuId int64) ([]model.ProductModel, error)
}

type StocksClient interface {
	StockInfo(ctx context.Context, skuId int64) (uint32, error)
}

type OrdersClient interface {
	CreateOrder(ctx context.Context, userId int64, items []model.CartItemModel) (int64, error)
}

type CartService struct {
	cartRepository CartRepository
	productService ProductService
	ordersClient   OrdersClient
	validator      *validator.Validate
}

func NewCartService(
	cartRepository CartRepository,
	productService ProductService,
	orders OrdersClient,
) *CartService {
	return &CartService{
		cartRepository: cartRepository,
		productService: productService,
		ordersClient:   orders,
		validator:      validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Create implements create_cart_item_handler.CartService.
func (service *CartService) Create(ctx context.Context, cartItem *model.CartItemModel) (bool, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "cart_service.Create")
	defer span.End()

	if err := service.validator.Struct(cartItem); err != nil {
		logger.Warn("Failed to validate cart item", "error", err)
		return false, fmt.Errorf("failed to validate cart item %w", err)
	}

	exists, err := service.productService.IsProductExists(ctx, cartItem.SkuId)
	if err != nil {
		logger.Warn("Failed to check if product exists", "error", err)
		return false, fmt.Errorf("failed to check if product exists %w", err)
	}

	if !exists {
		return false, model.ErrProductDoesNotExist
	}

	return service.cartRepository.CreateItem(ctx, cartItem)
}

// Checkout implements checkout_handler.CartService.
func (service *CartService) Checkout(ctx context.Context, userId int64) (int64, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "cart_service.Checkout")
	defer span.End()

	if userId <= 0 {
		logger.Warn("Checkout failed, userId less than 0", "userId", userId)
		return 0, fmt.Errorf("checkout: userId: %d, %w", userId, model.ErrorUserIdLessThanZero)
	}

	items := service.cartRepository.GetAllOrderBySku(ctx, userId)
	if len(items) == 0 {
		logger.Debug("Checkout failed, cart is empty", "userId", userId)
		return 0, fmt.Errorf("checkout: %w", model.ErrTheCartIsEmpty)
	}

	orderId, err := service.ordersClient.CreateOrder(ctx, userId, items)
	if err != nil {
		logger.Warn("Checkout failed, create order failed", "userId", userId, "error", err)
		return 0, fmt.Errorf("checkout: create order failed for user %d, %w", userId, err)
	}

	service.cartRepository.DeleteAll(ctx, userId)
	return orderId, nil
}

// DeleteAll implements delete_cart_handler.CartService.
func (service *CartService) DeleteAll(ctx context.Context, userId int64) error {
	ctx, span := otel.Tracer("service").Start(ctx, "cart_service.DeleteAll")
	defer span.End()

	if userId <= 0 {
		logger.Warn("DeleteAll failed, userId less than 0", "userId", userId)
		return fmt.Errorf("deleteAll: userId %d, %w", userId, model.ErrorUserIdLessThanZero)
	}

	service.cartRepository.DeleteAll(ctx, userId)
	return nil
}

// DeleteBySkuId implements delete_cart_item_handler.CartService.
func (service *CartService) DeleteBySkuId(ctx context.Context, userId int64, skuId int64) error {
	ctx, span := otel.Tracer("service").Start(ctx, "cart_service.DeleteBySkuId")
	defer span.End()

	if userId <= 0 {
		logger.Warn("DeleteBySkuId failed, userId less than 0", "userId", userId)
		return fmt.Errorf("deleteBySku: userId: %d, %w", userId, model.ErrorUserIdLessThanZero)
	}

	if skuId <= 0 {
		logger.Warn("DeleteBySkuId failed, skuId less than 0", "skuId", skuId)
		return fmt.Errorf("deleteBySku: sku: %d, %w", userId, model.ErrorSkuIdLessThanZero)
	}

	service.cartRepository.DeleteBySku(ctx, userId, skuId)
	return nil
}

func (service *CartService) GetAllOrderBySku(ctx context.Context, userId int64) (model.AllCartItemsModel, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "cart_service.GetAllOrderBySku")
	defer span.End()

	if userId <= 0 {
		logger.Warn("GetAllOrderBySku failed, userId less than 0", "userId", userId)
		return model.AllCartItemsModel{}, fmt.Errorf("getAllOrderBySku: userId %d, %w", userId, model.ErrorUserIdLessThanZero)
	}

	var items = service.cartRepository.GetAllOrderBySku(ctx, userId)
	if len(items) == 0 {
		logger.Debug("GetAllOrderBySku failed, cart is empty", "userId", userId)
		return model.AllCartItemsModel{}, fmt.Errorf("getAllOrderBySku: %w", model.ErrTheCartIsEmpty)
	}

	return service.enrichCartProducts(ctx, items)
}

// We need to speak with the products team to provide a better API, ugly n+1
func (service *CartService) enrichCartProducts(ctx context.Context, items []model.CartItemModel) (model.AllCartItemsModel, error) {
	var totalPrice atomic.Uint32
	var promiseGroup = route_err_group.NewRouteErrorGroup[model.EnrichedCartItemModel](ctx,
		route_err_group.Options{Rps: 10, BufferSize: len(items)})

	for cartItemIndex := range items {
		promiseGroup.Run(func(ctx context.Context) (model.EnrichedCartItemModel, error) {
			var cartProductSku = items[cartItemIndex].SkuId
			product, err := service.productService.GetProduct(ctx, cartProductSku)
			if err != nil {
				return model.EnrichedCartItemModel{}, fmt.Errorf("product not found: %w", err)
			}

			totalPrice.Add(product.Price * items[cartItemIndex].Count)
			return model.EnrichedCartItemModel{
				SkuId: items[cartItemIndex].SkuId,
				Count: items[cartItemIndex].Count,
				Name:  product.Name,
				Price: product.Price,
			}, nil
		})
	}

	allItems, err := promiseGroup.Await()
	if err != nil {
		return model.AllCartItemsModel{}, fmt.Errorf("error found while awaiting a products fetch, %w", err)
	}

	return model.AllCartItemsModel{
		Items: allItems,
		Total: totalPrice.Load(),
	}, nil
}

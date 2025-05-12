package service_test

import (
	"context"
	"errors"
	"reflect"
	"route256/cart/internal/domain/cart/service"
	"route256/cart/internal/domain/model"
	"testing"

	"github.com/gojuno/minimock/v3"
)

func TestCartService_Create(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	type fields struct {
		cartRepository service.CartRepository
		productService service.ProductService
		ordersClient   service.OrdersClient
	}
	tests := []struct {
		name    string
		fields  fields
		model   *model.CartItemModel
		want    bool
		wantErr bool
	}{
		{
			name: "should add new item and call repository",
			fields: fields{
				cartRepository: NewCartRepositoryMock(mc).
					CreateItemMock.Return(true, nil),
				productService: NewProductServiceMock(mc).
					IsProductExistsMock.Return(true, nil),
				ordersClient: NewOrdersClientMock(t),
			},
			model:   &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1},
			want:    true,
			wantErr: false,
		},
		{
			name: "should return error if model is not valid",
			fields: fields{
				cartRepository: NewCartRepositoryMock(mc),
				productService: NewProductServiceMock(mc),
				ordersClient:   NewOrdersClientMock(t),
			},
			model:   &model.CartItemModel{UserId: 1, SkuId: -1, Count: 1},
			want:    false,
			wantErr: true,
		},
		{
			name: "should return error if product does not exist",
			fields: fields{
				cartRepository: NewCartRepositoryMock(mc),
				productService: NewProductServiceMock(mc).
					IsProductExistsMock.Return(false, nil),
				ordersClient: NewOrdersClientMock(t),
			},
			model:   &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1},
			want:    false,
			wantErr: true,
		},
		{
			name: "should return error from product service",
			fields: fields{
				cartRepository: NewCartRepositoryMock(mc),
				productService: NewProductServiceMock(mc).
					IsProductExistsMock.Return(false, errors.New("some error")),
				ordersClient: NewOrdersClientMock(t),
			},
			model:   &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := service.NewCartService(tt.fields.cartRepository, tt.fields.productService,
				tt.fields.ordersClient)
			got, err := service.Create(context.Background(), tt.model)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.Create() = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CartService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCartService_DeleteAll(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	type fields struct {
		cartRepository service.CartRepository
		productService service.ProductService
		ordersClient   service.OrdersClient
	}
	tests := []struct {
		name    string
		fields  fields
		userId  int64
		wantErr bool
	}{
		{
			name:   "should remove item and call repository",
			userId: 1,
			fields: fields{
				cartRepository: NewCartRepositoryMock(mc).
					DeleteAllMock.Return(),
				productService: NewProductServiceMock(mc),
				ordersClient:   NewOrdersClientMock(t),
			},
			wantErr: false,
		},
		{
			name:   "should fail if user id is invalid",
			userId: -1,
			fields: fields{
				cartRepository: NewCartRepositoryMock(mc),
				productService: NewProductServiceMock(mc),
				ordersClient:   NewOrdersClientMock(t),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := service.NewCartService(tt.fields.cartRepository, tt.fields.productService,
				tt.fields.ordersClient)

			if err := service.DeleteAll(context.Background(), tt.userId); (err != nil) != tt.wantErr {
				t.Errorf("CartService.DeleteAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCartService_DeleteBySkuId(t *testing.T) {
	t.Parallel()
	type fields struct {
		cartRepository service.CartRepository
		productService service.ProductService
		ordersClient   service.OrdersClient
		stocksClient   service.StocksClient
	}
	type args struct {
		userId int64
		skuId  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should remove the item by sky",
			args: args{userId: 1, skuId: 1},
			fields: fields{
				cartRepository: NewCartRepositoryMock(t).
					DeleteBySkuMock.Return(),
				productService: NewProductServiceMock(t),
				ordersClient:   NewOrdersClientMock(t),
				stocksClient:   NewStocksClientMock(t),
			},
			wantErr: false,
		},
		{
			name: "should return error if userId is invalid",
			args: args{userId: 0, skuId: 1},
			fields: fields{
				cartRepository: NewCartRepositoryMock(t),
				productService: NewProductServiceMock(t),
				ordersClient:   NewOrdersClientMock(t),
				stocksClient:   NewStocksClientMock(t),
			},
			wantErr: true,
		},
		{
			name: "should return error if skuId is invalid",
			args: args{userId: 1, skuId: -1},
			fields: fields{
				cartRepository: NewCartRepositoryMock(t),
				productService: NewProductServiceMock(t),
				ordersClient:   NewOrdersClientMock(t),
				stocksClient:   NewStocksClientMock(t),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := service.NewCartService(tt.fields.cartRepository, tt.fields.productService,
				tt.fields.ordersClient)
			if err := service.DeleteBySkuId(context.Background(), tt.args.userId, tt.args.skuId); (err != nil) != tt.wantErr {
				t.Errorf("CartService.DeleteBySkuId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCartService_GetAllOrderBySku(t *testing.T) {
	t.Parallel()
	type fields struct {
		cartRepository service.CartRepository
		productService service.ProductService
		orderClient    service.OrdersClient
		stocksClient   service.StocksClient
	}
	tests := []struct {
		name    string
		fields  fields
		userId  int64
		want    model.AllCartItemsModel
		wantErr bool
	}{
		{
			name:   "should return all items from cart merged with products",
			userId: 1,
			fields: fields{
				cartRepository: NewCartRepositoryMock(t).
					GetAllOrderBySkuMock.Return([]model.CartItemModel{
					{UserId: 1, SkuId: 1, Count: 1},
					{UserId: 1, SkuId: 2, Count: 1},
					{UserId: 1, SkuId: 3, Count: 1},
				}),
				productService: NewProductServiceMock(t).
					GetProductMock.When(minimock.AnyContext, 1).Then(model.ProductModel{
					Name:  "product1",
					SkuId: 1,
					Price: 100,
				}, nil).GetProductMock.When(minimock.AnyContext, 2).Then(model.ProductModel{
					Name:  "product2",
					SkuId: 2,
					Price: 200,
				}, nil).GetProductMock.When(minimock.AnyContext, 3).Then(model.ProductModel{
					Name:  "product3",
					SkuId: 3,
					Price: 300,
				}, nil),
				orderClient:  NewOrdersClientMock(t),
				stocksClient: NewStocksClientMock(t),
			},
			want: model.AllCartItemsModel{
				Items: []model.EnrichedCartItemModel{
					{SkuId: 1, Count: 1, Name: "product1", Price: 100},
					{SkuId: 2, Count: 1, Name: "product2", Price: 200},
					{SkuId: 3, Count: 1, Name: "product3", Price: 300},
				},
				Total: 600,
			},
			wantErr: false,
		},
		{
			name:   "should return error if userId is invalid",
			userId: 0,
			fields: fields{
				cartRepository: NewCartRepositoryMock(t),
				productService: NewProductServiceMock(t),
				orderClient:    NewOrdersClientMock(t),
				stocksClient:   NewStocksClientMock(t),
			},
			want:    model.AllCartItemsModel{},
			wantErr: true,
		},
		{
			name:   "should return error if userId is invalid",
			userId: 1,
			fields: fields{
				cartRepository: NewCartRepositoryMock(t).
					GetAllOrderBySkuMock.Return([]model.CartItemModel{}),
				productService: NewProductServiceMock(t),
				orderClient:    NewOrdersClientMock(t),
				stocksClient:   NewStocksClientMock(t),
			},
			want:    model.AllCartItemsModel{},
			wantErr: true,
		},
		{
			name:   "should return error if userId is invalid",
			userId: 1,
			fields: fields{
				cartRepository: NewCartRepositoryMock(t).
					GetAllOrderBySkuMock.Return([]model.CartItemModel{
					{UserId: 1, SkuId: 1, Count: 1},
				}),
				productService: NewProductServiceMock(t).
					GetProductMock.Return(model.ProductModel{}, errors.New("get product error")),
				orderClient:  NewOrdersClientMock(t),
				stocksClient: NewStocksClientMock(t),
			},
			want:    model.AllCartItemsModel{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := service.NewCartService(tt.fields.cartRepository, tt.fields.productService,
				tt.fields.orderClient)

			got, err := service.GetAllOrderBySku(context.Background(), tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.GetAllOrderBySku() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartService.GetAllOrderBySku() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCartService_Checkout(t *testing.T) {
	t.Parallel()
	type fields struct {
		cartRepository service.CartRepository
		productService service.ProductService
		orderClient    service.OrdersClient
	}
	tests := []struct {
		name    string
		fields  fields
		userId  int64
		get     int64
		wantErr bool
	}{
		{
			name:   "should successfully checkout by userId",
			userId: 1,
			fields: fields{
				cartRepository: NewCartRepositoryMock(t).
					GetAllOrderBySkuMock.Return([]model.CartItemModel{
					{UserId: 1, SkuId: 1, Count: 1}}).
					DeleteAllMock.Return(),
				productService: NewProductServiceMock(t),
				orderClient: NewOrdersClientMock(t).
					CreateOrderMock.When(minimock.AnyContext, 1, []model.CartItemModel{
					{UserId: 1, SkuId: 1, Count: 1}}).Then(321, nil),
			},
			get:     321,
			wantErr: false,
		},
		{
			name:   "should return error if userId is invalid",
			userId: 0,
			fields: fields{
				cartRepository: NewCartRepositoryMock(t),
				productService: NewProductServiceMock(t),
				orderClient:    NewOrdersClientMock(t),
			},
			get:     0,
			wantErr: true,
		},
		{
			name:   "should return error if failed to create order",
			userId: 1,
			fields: fields{
				cartRepository: NewCartRepositoryMock(t).
					GetAllOrderBySkuMock.Return([]model.CartItemModel{
					{UserId: 1, SkuId: 1, Count: 1}}),
				productService: NewProductServiceMock(t),
				orderClient: NewOrdersClientMock(t).
					CreateOrderMock.Return(0, errors.New("create order error")),
			},
			get:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := service.NewCartService(tt.fields.cartRepository, tt.fields.productService,
				tt.fields.orderClient)

			orderId, err := service.Checkout(context.Background(), tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.Checkout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if orderId != tt.get {
				t.Errorf("CartService.Checkout() = %v, want %v", orderId, tt.get)
			}
		})
	}
}

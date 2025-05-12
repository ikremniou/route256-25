package repository_test

import (
	"context"
	"math"
	"reflect"
	"route256/cart/internal/domain/cart/repository"
	"route256/cart/internal/domain/model"
	"testing"
)

func TestCartRepository_CreateItem(t *testing.T) {
	t.Parallel()
	type args struct {
		item *model.CartItemModel
	}
	type call_attempt struct {
		args    args
		want    bool
		wantErr bool
	}
	tests := []struct {
		name  string
		tries []call_attempt
	}{
		{
			name: "should add new item and return true",
			tries: []call_attempt{
				{
					args:    args{item: &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1}},
					want:    true,
					wantErr: false,
				},
			},
		},
		{
			name: "should add existing item and return false",
			tries: []call_attempt{
				{
					args:    args{item: &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1}},
					want:    true,
					wantErr: false,
				},
				{
					args:    args{item: &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1}},
					want:    false,
					wantErr: false,
				},
			},
		},
		{
			name: "should add 2 items to the users cart and return true",
			tries: []call_attempt{
				{
					args:    args{item: &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1}},
					want:    true,
					wantErr: false,
				},
				{
					args:    args{item: &model.CartItemModel{UserId: 1, SkuId: 2, Count: 1}},
					want:    true,
					wantErr: false,
				},
			},
		},
		{
			name: "should return error in case of the count overflow",
			tries: []call_attempt{
				{
					args:    args{item: &model.CartItemModel{UserId: 1, SkuId: 1, Count: math.MaxUint32}},
					want:    true,
					wantErr: false,
				},
				{
					args:    args{item: &model.CartItemModel{UserId: 1, SkuId: 1, Count: math.MaxUint32}},
					want:    false,
					wantErr: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := repository.NewCartRepository()

			for _, attempt := range tt.tries {
				got, err := c.CreateItem(context.Background(), attempt.args.item)

				if (err != nil) != attempt.wantErr {
					t.Errorf("CartRepository.CreateItem() error = %v, wantErr %v", err, attempt.wantErr)
					return
				}
				if got != attempt.want {
					t.Errorf("CartRepository.CreateItem() = %v, want %v", got, attempt.want)
				}
			}
		})
	}
}

func TestCartRepository_DeleteBySku(t *testing.T) {
	t.Parallel()
	type args struct {
		userId int64
		skuId  int64
	}
	tests := []struct {
		name   string
		userId int64
		having []model.CartItemModel
		tries  []args
		has    []model.CartItemModel
	}{
		{
			name:   "should remove single user cart item",
			userId: 1,
			having: []model.CartItemModel{
				{UserId: 1, SkuId: 1, Count: 1},
			},
			tries: []args{{userId: 1, skuId: 1}},
			has:   []model.CartItemModel{},
		},
		{
			name:   "should remove single user cart if user has 2 items",
			userId: 1,
			having: []model.CartItemModel{
				{UserId: 1, SkuId: 1, Count: 1},
				{UserId: 1, SkuId: 2, Count: 1},
			},
			tries: []args{{userId: 1, skuId: 1}},
			has: []model.CartItemModel{
				{UserId: 1, SkuId: 2, Count: 1},
			},
		},
		{
			name:   "should do nothing if user has no items",
			userId: 1,
			having: []model.CartItemModel{},
			tries:  []args{{userId: 1, skuId: 1}},
			has:    []model.CartItemModel{},
		},
		{
			name:   "should remove single user cart if user has 2 items",
			userId: 1,
			having: []model.CartItemModel{
				{UserId: 1, SkuId: 1, Count: 1},
				{UserId: 1, SkuId: 2, Count: 1},
			},
			tries: []args{{userId: 1, skuId: 1}},
			has: []model.CartItemModel{
				{UserId: 1, SkuId: 2, Count: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := repository.NewCartRepository()
			for _, item := range tt.having {
				_, err := c.CreateItem(context.Background(), &item)
				if err != nil {
					t.Errorf("CartRepository.CreateItem() error = %v", err)
				}
			}

			for _, attempt := range tt.tries {
				c.DeleteBySku(context.Background(), attempt.userId, attempt.skuId)
			}

			if items := c.GetAllOrderBySku(context.Background(), tt.userId); !compareCartItems(items, tt.has) {
				t.Errorf("CartRepository.GetAllOrderBySku() = %v, want %v", items, tt.has)
			}
		})
	}
}

func TestCartRepository_DeleteAll(t *testing.T) {
	t.Parallel()
	type args struct {
		userId int64
	}
	tests := []struct {
		name   string
		having []model.CartItemModel
		args   args
	}{
		{
			name: "should remove the all items from the user cart",
			having: []model.CartItemModel{
				{UserId: 1, SkuId: 1, Count: 1},
				{UserId: 1, SkuId: 2, Count: 1},
				{UserId: 1, SkuId: 3, Count: 1},
			},
			args: args{userId: 1},
		},
		{
			name:   "should do nothing if user has no items",
			having: []model.CartItemModel{},
			args:   args{userId: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := repository.NewCartRepository()
			for _, item := range tt.having {
				_, err := c.CreateItem(context.Background(), &item)
				if err != nil {
					t.Errorf("CartRepository.CreateItem() error = %v", err)
				}
			}

			c.DeleteAll(context.Background(), tt.args.userId)
			items := c.GetAllOrderBySku(context.Background(), tt.args.userId)

			if len(items) != 0 {
				t.Errorf("CartRepository.GetAllOrderBySku() = having %v, want 0", items)
			}
		})
	}
}

func TestCartRepository_GetAllOrderBySku(t *testing.T) {
	t.Parallel()
	type args struct {
		userId int64
	}
	tests := []struct {
		name   string
		args   args
		having []model.CartItemModel
		want   []model.CartItemModel
	}{
		{
			name: "should return a single item from cart",
			args: args{userId: 1},
			having: []model.CartItemModel{
				{UserId: 1, SkuId: 1, Count: 1},
			},
			want: []model.CartItemModel{
				{UserId: 1, SkuId: 1, Count: 1},
			},
		},
		{
			name:   "should return empty array if user has no items",
			args:   args{userId: 1},
			having: []model.CartItemModel{},
			want:   []model.CartItemModel{},
		},
		{
			name: "should return items in sort order",
			args: args{userId: 1},
			having: []model.CartItemModel{
				{UserId: 1, SkuId: 99, Count: 1},
				{UserId: 1, SkuId: 9, Count: 1},
				{UserId: 1, SkuId: 1, Count: 1},
			},
			want: []model.CartItemModel{
				{UserId: 1, SkuId: 1, Count: 1},
				{UserId: 1, SkuId: 9, Count: 1},
				{UserId: 1, SkuId: 99, Count: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := repository.NewCartRepository()

			for _, item := range tt.having {
				_, err := c.CreateItem(context.Background(), &item)
				if err != nil {
					t.Errorf("CartRepository.CreateItem() error = %v", err)
				}
			}

			if got := c.GetAllOrderBySku(context.Background(), tt.args.userId); !compareCartItems(got, tt.want) {
				t.Errorf("CartRepository.GetAllOrderBySku() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareCartItems(a, b []model.CartItemModel) bool {
	if len(a) != len(b) {
		return false
	}

	if len(a) == 0 && len(b) == 0 {
		return true
	}

	return reflect.DeepEqual(a, b)
}

package products_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/domain/products"
	"testing"

	"github.com/gojuno/minimock/v3"
)

func TestProductsClient_GetProductsAot(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	type args struct {
		count      int64
		startSkuId int64
	}
	tests := []struct {
		name      string
		args      args
		transport http.RoundTripper
		want      []model.ProductModel
		wantErr   bool
	}{
		{
			name: "should return array of products",
			args: args{count: 2, startSkuId: 1},
			want: []model.ProductModel{
				{SkuId: 101, Name: "name1", Price: 100},
				{SkuId: 102, Name: "name2", Price: 200},
			},
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader(
					[]byte(`[{"sku":101,"name":"name1","price":100},{"sku":102,"name":"name2","price":200}]`))),
			}, nil),
			wantErr: false,
		},
		{
			name: "should fail on the invalid json",
			args: args{count: 2, startSkuId: 1},
			want: []model.ProductModel{},
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader(
					[]byte(`///`))),
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := *products.NewProductsClientForTest(tt.transport, "test", "test")
			got, err := client.GetProductsAot(context.Background(), tt.args.count, tt.args.startSkuId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductsClient.GetProductsAot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductsClient.GetProductsAot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductsClient_GetProduct(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	tests := []struct {
		name      string
		skuId     int64
		transport http.RoundTripper
		want      model.ProductModel
		wantErr   bool
	}{
		{
			name:  "should return product by sku id",
			skuId: 1,
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader(
					[]byte(`{"sku":1,"name":"name1","price":100}`))),
			}, nil),
			want:    model.ProductModel{SkuId: 1, Name: "name1", Price: 100},
			wantErr: false,
		},
		{
			name:  "should return error on the invalid json",
			skuId: 1,
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader(
					[]byte(`///`))),
			}, nil),
			want:    model.ProductModel{},
			wantErr: true,
		},
		{
			name:  "should return error when request failed",
			skuId: 1,
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{}, errors.New("error")),
			want:    model.ProductModel{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := products.NewProductsClientForTest(tt.transport, "test", "test")

			got, err := client.GetProduct(context.Background(), tt.skuId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductsClient.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductsClient.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductsClient_IsProductExists(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	tests := []struct {
		name      string
		skuId     int64
		transport http.RoundTripper
		want      bool
		wantErr   bool
	}{
		{
			name:  "should return 200 and true if product exists",
			skuId: 1,
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{StatusCode: http.StatusOK}, nil),
			want:    true,
			wantErr: false,
		},
		{
			name:  "should return false for non-200 status code",
			skuId: 1,
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{StatusCode: http.StatusNotFound}, nil),
			want:    false,
			wantErr: false,
		},
		{
			name:  "should return err if request failed",
			skuId: 1,
			transport: NewRoundTripperMock(mc).
				RoundTripMock.Return(&http.Response{}, errors.New("test")),
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := products.NewProductsClientForTest(tt.transport, "test", "test")

			got, err := client.IsProductExists(context.Background(), tt.skuId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductsClient.IsProductExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProductsClient.IsProductExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

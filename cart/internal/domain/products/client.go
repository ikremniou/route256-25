package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/cart_config"
	"route256/cart/internal/infra/sre"
	"route256/cart/internal/infra/tripper"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type ProductsClient struct {
	http    http.Client
	address string
	apiKey  string
}

const TwitterStatusCodeRateLimit = 420

func NewProductsClient(cartConfig *cart_config.Config) *ProductsClient {
	transport := http.DefaultTransport
	transport = tripper.NewRetryRoundTripper(transport, tripper.RetryConfig{
		RetryOn:   []int{http.StatusTooManyRequests, TwitterStatusCodeRateLimit},
		Times:     3,
		WaitForMs: 500,
	})
	address := fmt.Sprintf("http://%s:%s", cartConfig.Products.Host, cartConfig.Products.Port)
	transport = otelhttp.NewTransport(transport)

	return &ProductsClient{
		http:    http.Client{Transport: transport},
		address: address,
		apiKey:  cartConfig.Products.Token,
	}
}

// GetProductsAot implements service.ProductService.
func (client *ProductsClient) GetProductsAot(ctx context.Context, count int64, startSkuId int64) ([]model.ProductModel, error) {
	ctx, span := otel.Tracer("client").Start(ctx, "products_client.GetProductsAot")
	defer span.End()

	var url = fmt.Sprintf("%s/product?count=%d&start_after_sku=%d", client.address, count, startSkuId)
	startTime := time.Now()
	response, err := client.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	sre.TrackExternalRequest("products_get_products_aot", err, startTime)

	var products []model.ProductModel
	if err := json.NewDecoder(response.Body).Decode(&products); err != nil {
		return nil, err
	}

	return products, nil
}

func (client *ProductsClient) GetProduct(ctx context.Context, skuId int64) (model.ProductModel, error) {
	ctx, span := otel.Tracer("client").Start(ctx, "products_client.GetProduct")
	defer span.End()

	var url = fmt.Sprintf("%s/product/%d", client.address, skuId)
	startTime := time.Now()
	response, err := client.doRequest(ctx, url)
	if err != nil {
		return model.ProductModel{}, err
	}
	defer response.Body.Close()
	sre.TrackExternalRequest("products_get_product", err, startTime)

	if response.StatusCode != http.StatusOK {
		return model.ProductModel{}, fmt.Errorf("failed to get product %d", response.StatusCode)
	}

	var product model.ProductModel
	if err := json.NewDecoder(response.Body).Decode(&product); err != nil {
		return model.ProductModel{}, err
	}

	return product, nil
}

// IsProductExists implements service.ProductService.
func (client *ProductsClient) IsProductExists(ctx context.Context, skuId int64) (bool, error) {
	ctx, span := otel.Tracer("client").Start(ctx, "products_client.IsProductExists")
	defer span.End()

	var url = fmt.Sprintf("%s/product/%d", client.address, skuId)
	startTime := time.Now()
	response, err := client.doRequest(ctx, url)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()
	sre.TrackExternalRequest("products_is_product_exists", err, startTime)
	return response.StatusCode == http.StatusOK, nil
}

func (client *ProductsClient) doRequest(ctx context.Context, url string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-API-KEY", client.apiKey)
	response, err := client.http.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

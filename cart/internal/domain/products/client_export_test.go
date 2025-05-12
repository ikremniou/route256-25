package products

import (
	"net/http"
)

func NewProductsClientForTest(transport http.RoundTripper, address, apiKey string) *ProductsClient {
	return &ProductsClient{
		http:    http.Client{Transport: transport},
		address: address,
		apiKey:  apiKey,
	}
}

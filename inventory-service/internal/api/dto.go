package api

import (
	"time"

	"github.com/Daty26/order-system/inventory-service/internal/model"
)

type InsertProductRequest struct {
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	PriceCents int64  `json:"price_cents"`
}
type ProductResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToProductResponse(product model.Product) ProductResponse {
	return ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Quantity:   product.Quantity,
		PriceCents: product.PriceCents,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}
}
func ToProductResponses(products []model.Product) []ProductResponse {
	productResponses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}
	return productResponses
}

type QuoteProductsRequest struct {
	IDs []int `json:"ids"`
}
type QuoteProductsResponse struct {
	IDs        int   `json:"id"`
	PriceCents int64 `json:"price_cents"`
}

func ToQuatesProductResponse(productQuotes model.ProductQuote) QuoteProductsResponse {
	return QuoteProductsResponse{
		IDs:        productQuotes.ID,
		PriceCents: productQuotes.PriceCents,
	}
}
func ToQuotesProductReponses(productQuotesModels []model.ProductQuote) []QuoteProductsResponse {
	productQuotesResponses := make([]QuoteProductsResponse, 0, len(productQuotesModels))
	for _, productQuotesModel := range productQuotesModels {
		productQuotesResponses = append(productQuotesResponses, ToQuatesProductResponse(productQuotesModel))
	}
	return productQuotesResponses
}

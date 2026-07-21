package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Daty26/order-system/order-service/internal/service"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), httpClient: httpClient}
}

type productsResponse struct {
	Data  []quoteResponse `json:"data"`
	Error string          `json:"error"`
}

type quoteResponse struct {
	ID         int   `json:"id"`
	PriceCents int64 `json:"price_cents"`
}
type quoteProductsRequest struct {
	IDs []int `json:"ids"`
}

func (c *Client) GetQuotes(ctx context.Context, ids []int) (map[int]service.Product, error) {
	body := quoteProductsRequest{
		IDs: ids,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal quote request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/products/quotes",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("create inventory req: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call inventory service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("inventory service returned status: %d", resp.StatusCode)
	}
	var decoded productsResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode inventory service: %w", err)
	}

	quotes := make(map[int]service.Product, len(ids))
	for _, product := range decoded.Data {
		quotes[product.ID] = service.Product{
			ID:         product.ID,
			PriceCents: product.PriceCents,
		}
	}

	return quotes, nil
}

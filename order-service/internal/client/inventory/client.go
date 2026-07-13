package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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
	Data  []productResponse `json:"data"`
	Error string            `json:"error"`
}

type productResponse struct {
	ID    int     `json:"id"`
	Price float64 `json:"price"`
}

func (c *Client) GetQuotes(ctx context.Context, ids []int) (map[int]service.Product, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/products", nil)
	if err != nil {
		return nil, fmt.Errorf("create inventory req: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call inventory service: %w", err)
	}
	defer req.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("inventory service returned status: %d", resp.StatusCode)
	}
	var body productsResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode inventory service: %w", err)
	}
	wanted := make(map[int]struct{}, len(ids))

	for _, id := range ids {
		wanted[id] = struct{}{}
	}
	quotes := make(map[int]service.Product, len(ids))
	for _, product := range body.Data {
		if _, ok := wanted[product.ID]; !ok {
			continue
		}
		quotes[product.ID] = service.Product{
			ID:         product.ID,
			PriceCents: int64(math.Round(product.Price * 100)),
		}
	}
	return quotes, nil
}

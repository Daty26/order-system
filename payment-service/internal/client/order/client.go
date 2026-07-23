package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Daty26/order-system/payment-service/internal/service"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, client *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: client,
	}
}

type OrderIDRequest struct {
	ID int `json:"id"`
}
type orderResponse struct {
	Data  service.OrderSummary `json:"data"`
	Error string               `json:"error"`
}

func (c *Client) GetOrder(ctx context.Context, orderID int, authHeader string) (service.OrderSummary, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/orders/"+strconv.Itoa(orderID),
		nil,
	)
	if err != nil {
		return service.OrderSummary{}, fmt.Errorf("create get order req: %w", err)
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return service.OrderSummary{}, fmt.Errorf("call order service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return service.OrderSummary{}, fmt.Errorf("order service returned status: %d", resp.StatusCode)
	}
	var decoded orderResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return service.OrderSummary{}, fmt.Errorf("decode order service: %w", err)
	}
	return decoded.Data, nil

}

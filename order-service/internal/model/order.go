package model

import "time"

type Orders struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Status    string      `json:"string"`
	CreatedAt time.Time   `json:"created_at"`
	Items     []OrderItem `json:"items"`
}
type OrderItem struct {
	ID        int `json:"id"`
	OrderId   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

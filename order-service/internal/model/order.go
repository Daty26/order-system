package model

import "time"

type Orders struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	Items     []OrderItems `json:"items"`
}
type OrderItems struct {
	ID        int `json:"id"`
	OrderId   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

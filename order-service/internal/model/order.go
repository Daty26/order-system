package model

type Order struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"amount"`
	UserID    int `json:"user_id"`
}

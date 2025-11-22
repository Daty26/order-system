package model

type Order struct {
	ID       int    `json:"id"`
	Item     string `json:"item"`
	Quantity int    `json:"amount"`
	UserID   int    `json:"user_id"`
}

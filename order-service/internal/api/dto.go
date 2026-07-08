package api

type CreatedOrderRequest struct {
	Items []CreatedOrderItemRequest `json:"items"`
}

type CreatedOrderItemRequest struct{
	ProductID int `json:"product_id"`
	Quantity int `json:"quantity"`
}

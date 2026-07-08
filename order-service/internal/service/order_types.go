package service

import "context"

type CreatedOrderInput struct{
	UserID int
	Items []CreatedOrderItemInput
}
type CreatedOrderItemInput struct{
	ProductID int
	Quantity int
}
type Product struct{
	ID int
	PriceCents int64
}
type ProductInventory interface{
	GetQuotes(ctx context.Context, ids []int) (map[int]Product, error)
}
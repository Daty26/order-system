package service

type InsertProductInput struct {
	Name       string
	Quantity   int
	PriceCents int64
}
type UpdateProductInput struct {
	ID         int
	PriceCents int64
}
type UpdateQuantityInput struct {
	ID       int
	Quantity int
}
type GetQuotesInput struct {
	IDs []int
}

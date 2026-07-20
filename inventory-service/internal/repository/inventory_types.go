package repository

type InsertProductParams struct {
	Name       string
	Quantity   int
	PriceCents int64
}
type UpdateQuantityParams struct {
	ProductID int
	Quantity  int
}
type UpdatePriceCentsParams struct {
	ProductID  int
	PriceCents int64
}
type ReduceStockParams struct {
	ProductID int
	Quantity  int
}

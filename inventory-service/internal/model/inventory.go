package model

import "time"

type Product struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
type ProductQuote struct {
	ID         int
	PriceCents int64
}

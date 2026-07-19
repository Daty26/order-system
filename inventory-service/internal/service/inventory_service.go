package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daty26/order-system/inventory-service/internal/model"
	"github.com/Daty26/order-system/inventory-service/internal/repository"
)

type InventoryService struct {
	repo *repository.PostgresInventoryRepo
}

func NewInventoryService(inventoryRepo *repository.PostgresInventoryRepo) *InventoryService {
	return &InventoryService{repo: inventoryRepo}
}

func (is *InventoryService) GetAll(ctx context.Context, limit, offset int) ([]model.Product, error) {
	return is.repo.GetAll(ctx, limit, offset)
}

func (is *InventoryService) UpdateQuantity(ctx context.Context, id int, quantity int) (model.Product, error) {
	if quantity < 0 {
		return model.Product{}, ErrInvalidInput
	}
	return is.repo.UpdateQuantity(ctx, id, quantity)
}
func (s *InventoryService) InsertProduct(ctx context.Context, input InsertProductInput) (model.Product, error) {
	if input.Name == "" {
		return model.Product{}, ErrInvalidInput
	}
	if input.Quantity < 0 || input.PriceCents < 0 {
		return model.Product{}, ErrInvalidInput
	}
	params := repository.InsertProductParams{
		Name:       input.Name,
		Quantity:   input.Quantity,
		PriceCents: input.PriceCents,
	}
	return s.repo.Insert(ctx, params)
}
func (is *InventoryService) UpdatePrice(ctx context.Context, id int, price int64) (model.Product, error) {
	if id < 0 {
		return model.Product{}, ErrInvalidInput
	}
	if price < 0 {
		return model.Product{}, ErrInvalidInput
	}
	return is.repo.UpdatePriceCents(ctx, id, price)
}

func (s *InventoryService) ReduceStock(ctx context.Context, productId, quantity int) (model.Product, error) {
	if productId <= 0 || quantity <= 0 {
		return model.Product{}, ErrInvalidInput
	}
	updatedProduct, err := s.repo.ReduceStock(ctx, productId, quantity)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Product{}, ErrInsufficientStock
	}
	if err != nil {
		return model.Product{}, fmt.Errorf("reduce stock: %w", err)
	}
	return updatedProduct, nil
}

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
	repo repository.InventoryRepository
}

func NewInventoryService(inventoryRepo *repository.PostgresInventoryRepo) *InventoryService {
	return &InventoryService{repo: inventoryRepo}
}

func (is *InventoryService) GetAll(ctx context.Context, limit, offset int) ([]model.Product, error) {
	return is.repo.GetAll(ctx, limit, offset)
}

func (is *InventoryService) UpdateQuantity(ctx context.Context, input UpdateQuantityInput) (model.Product, error) {
	if input.ID < 0 || input.Quantity < 0 {
		return model.Product{}, ErrInvalidInput
	}

	params := repository.UpdateQuantityParams{
		ProductID: input.ID,
		Quantity:  input.Quantity,
	}

	return is.repo.UpdateQuantity(ctx, params)
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
func (is *InventoryService) UpdatePrice(ctx context.Context, input UpdateProductInput) (model.Product, error) {
	if input.ID < 0 || input.PriceCents < 0 {
		return model.Product{}, ErrInvalidInput
	}
	params := repository.UpdatePriceCentsParams{
		ProductID:  input.ID,
		PriceCents: input.PriceCents,
	}
	return is.repo.UpdatePriceCents(ctx, params)
}

func (s *InventoryService) ReduceStock(ctx context.Context, productId, quantity int) (model.Product, error) {
	if productId <= 0 || quantity <= 0 {
		return model.Product{}, ErrInvalidInput
	}
	params := repository.ReduceStockParams{
		ProductID: productId,
		Quantity:  quantity,
	}
	updatedProduct, err := s.repo.ReduceStock(ctx, params)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Product{}, ErrInsufficientStock
	}
	if err != nil {
		return model.Product{}, fmt.Errorf("reduce stock: %w", err)
	}
	return updatedProduct, nil
}

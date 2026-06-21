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

func (is *InventoryService) GetAll() ([]model.Product, error) {
	return is.repo.GetAll()
}
func (is *InventoryService) Insert(product model.Product) (model.Product, error) {
	if product.Quantity < 0 || product.Price < 0 {
		return model.Product{}, ErrInvalidInput
	}
	if len(product.Name) == 0 {
		return model.Product{}, errors.New("name can't be empty")
	}
	return is.repo.Insert(product)
}
func (is *InventoryService) UpdateQuantity(id int, quantity int) (model.Product, error) {
	if quantity < 0 {
		return model.Product{}, ErrInvalidInput
	}
	return is.repo.UpdateQuantity(id, quantity)
}
func (is *InventoryService) UpdatePrice(id int, price float64) (model.Product, error) {
	if id < 0 {
		return model.Product{}, errors.New("incorrect id")
	}
	if price < 0.0 {
		return model.Product{}, errors.New("incorrect price")
	}
	return is.repo.UpdatePrice(id, price)
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

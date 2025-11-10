package service

import (
	"errors"
	"github.com/Daty26/order-system/inventory-service/internal/model"
	"github.com/Daty26/order-system/inventory-service/internal/repository"
)

type InventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(inventoryService repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: inventoryService}
}

func (is *InventoryService) GetAll() ([]model.Product, error) {
	return is.repo.GetAll()
}
func (is *InventoryService) Insert(product model.Product) (model.Product, error) {
	if product.Quantity < 0 {
		return model.Product{}, errors.New("quantity can't be less than 0")
	}
	if len(product.Name) == 0 {
		return model.Product{}, errors.New("Name can't be empty")
	}
	return is.repo.Insert(product)
}
func (is *InventoryService) UpdateQuantity(id int, quantity int) (model.Product, error) {
	if id < 0 {
		return model.Product{}, errors.New("incorrect id")
	}
	if quantity < 0 {
		return model.Product{}, errors.New("quantity can't be less than 0")
	}
	return is.repo.UpdateQuantity(id, quantity)

}
func (is *InventoryService) ReduceStock(productId int, quantity int) error {
	product, err := is.repo.GetByID(productId)
	if err != nil {
		return err
	}
	if product.Quantity < quantity {
		return errors.New("the product is out of stock")
	}
	reducedQuantity := product.Quantity - quantity
	_, err = is.repo.UpdateQuantity(productId, reducedQuantity)
	if err != nil {
		return err
	}
	return nil
}

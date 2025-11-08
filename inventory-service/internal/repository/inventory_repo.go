package repository

import (
	"database/sql"
	"github.com/Daty26/order-system/inventory-service/internal/model"
)

type InventoryRepository interface {
	GetAll() ([]model.Product, error)
	Insert(product model.Product) (model.Product, error)
	UpdateQuantity(id int, quanity int) (model.Product, error)
}

type PostgresInventoryRepo struct {
	db *sql.DB
}

func NewPostgresInventoryRepo(db *sql.DB) *PostgresInventoryRepo {
	return &PostgresInventoryRepo{db: db}
}

func (pr *PostgresInventoryRepo) GetAll() ([]model.Product, error) {
	query := `select id, name, quantity, created_at, updated_at from inventory`
	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err = rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return []model.Product{}, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil
}
func (pr *PostgresInventoryRepo) Insert(product model.Product) (model.Product, error) {
	var insertedProduct model.Product
	query := `Insert into inventory (name, quantity) VALUES ($1, $2) RETURNING id, name, quantity, created_at, updated_at`
	err := pr.db.QueryRow(query, product.Name, product.Quantity).
		Scan(&insertedProduct.ID, &insertedProduct.Name, &insertedProduct.Quantity, &insertedProduct.CreatedAt, &insertedProduct.UpdatedAt)
	if err != nil {
		return model.Product{}, err
	}
	return insertedProduct, nil
}
func (pr *PostgresInventoryRepo) UpdateQuantity(id int, quantity int) (model.Product, error) {
	var updatedProduct model.Product
	query := `update inventory set quantity=$1 where id = $2 RETURNING id, name, quantity, created_at, updated_at`
	if err := pr.db.QueryRow(query, quantity, id).Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Quantity, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt); err != nil {
		return updatedProduct, err
	}
	return updatedProduct, nil

}

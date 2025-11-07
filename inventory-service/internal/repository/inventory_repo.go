package repository

import (
	"database/sql"
	"github.com/Daty26/order-system/inventory-service/internal/model"
)

type InventoryRepository interface {
	GetAll() ([]model.Product, error)
	Insert(product model.Product) (model.Product, error)
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
	query := `Insert into inventory (id, name, quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, quantity, created_at, updated_at`
	err := pr.db.QueryRow(query, product.ID, product.Name, product.Quantity, product.CreatedAt, product.UpdatedAt).
		Scan(&insertedProduct.ID, &insertedProduct.Name, &insertedProduct.Quantity, &insertedProduct.CreatedAt, &insertedProduct.UpdatedAt)
	if err != nil {
		return model.Product{}, err
	}
	return insertedProduct, nil

}

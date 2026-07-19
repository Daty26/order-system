package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daty26/order-system/inventory-service/internal/model"
)

type InventoryRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]model.Product, error)
	Insert(ctx context.Context, product model.Product) (model.Product, error)
	UpdateQuantity(ctx context.Context, id int, quanity int) (model.Product, error)
	GetByID(ctx context.Context, id int) (model.Product, error)
	UpdatePriceCents(ctx context.Context, id, priceCents int64) (model.Product, error)
	ReduceStock(ctx context.Context, id, quantity int) (model.Product, error)
}

type PostgresInventoryRepo struct {
	db *sql.DB
}

func NewPostgresInventoryRepo(db *sql.DB) *PostgresInventoryRepo {
	return &PostgresInventoryRepo{db: db}
}

func (pr *PostgresInventoryRepo) GetAll(ctx context.Context, limit, offset int) ([]model.Product, error) {
	query := `
		select id, name, quantity, price_cents, created_at, updated_at
		from inventory
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
`
	rows, err := pr.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err = rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.PriceCents, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return []model.Product{}, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil
}

func (pr *PostgresInventoryRepo) GetByID(ctx context.Context, id int) (model.Product, error) {
	var product model.Product
	query := `Select id, name, quantity, price_cents, created_at, updated_at from inventory where id=$1`
	if err := pr.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Quantity, &product.PriceCents, &product.CreatedAt, &product.UpdatedAt); err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (pr *PostgresInventoryRepo) Insert(ctx context.Context, product model.Product) (model.Product, error) {
	var insertedProduct model.Product
	query := `Insert into inventory (name, quantity, price_cents) VALUES ($1, $2, $3) RETURNING id, name, quantity, price_cents, created_at, updated_at`
	err := pr.db.QueryRowContext(ctx, query, product.Name, product.Quantity, product.PriceCents).
		Scan(&insertedProduct.ID, &insertedProduct.Name, &insertedProduct.Quantity, &insertedProduct.PriceCents, &insertedProduct.CreatedAt, &insertedProduct.UpdatedAt)
	if err != nil {
		return model.Product{}, err
	}
	return insertedProduct, nil
}

func (pr *PostgresInventoryRepo) UpdateQuantity(ctx context.Context, id int, quantity int) (model.Product, error) {
	var updatedProduct model.Product
	query := `
		update inventory
		set quantity=$1
		where id = $2
		RETURNING id, name, quantity, price_cents, created_at, updated_at`
	if err := pr.db.QueryRowContext(ctx, query, quantity, id).Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Quantity, &updatedProduct.PriceCents, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt); err != nil {
		return updatedProduct, err
	}
	return updatedProduct, nil
}

func (pr *PostgresInventoryRepo) UpdatePriceCents(ctx context.Context, id int, priceCents int64) (model.Product, error) {
	var updatedProduct model.Product
	query := `
		update inventory
		set price_cents=$1
		where id = $2
		RETURNING id, name, quantity, price_cents, created_at, updated_at
`
	if err := pr.db.QueryRowContext(ctx, query, priceCents, id).Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Quantity, &updatedProduct.PriceCents, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt); err != nil {
		return model.Product{}, err
	}
	return updatedProduct, nil
}

func (r *PostgresInventoryRepo) ReduceStock(ctx context.Context, id, quantity int) (model.Product, error) {
	var product model.Product
	const query = `
		UPDATE inventory
		SET quantity = quantity - $2,
			updated_at = now()
		WHERE id = $1 and quantity >= $2
		RETURNING id, name, quantity, price_cents, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, id, quantity).Scan(
		&product.ID,
		&product.Name,
		&product.Quantity,
		&product.PriceCents,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return model.Product{}, fmt.Errorf("reduce product stock: %w", err)
	}
	return product, nil
}

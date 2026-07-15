package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/Daty26/order-system/order-service/internal/db"
	"github.com/Daty26/order-system/order-service/internal/model"
	"github.com/lib/pq"
)

type OrderRep interface {
	Create(ctx context.Context, order model.Orders) (model.Orders, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Orders, error)
	GetByID(ctx context.Context, id int) (model.Orders, error)
	// Update(ctx context.Context, order model.Orders) (model.Orders, error)
	Delete(ctx context.Context, id int) error
	GetAllByUserID(ctx context.Context, userId, limit, offset int) ([]model.Orders, error)
}
type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}

func (r *PostgresOrderRepo) Create(ctx context.Context, order model.Orders) (model.Orders, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Orders{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const insertOrder = `
		insert into orders ( user_id, status, total_amount_cents) 
		values ($1, $2, $3) 
		RETURNING id, created_at
	`
	err = tx.QueryRowContext(ctx, insertOrder, order.UserID, order.Status, order.TotalAmountCents).Scan(&order.OrderID, &order.CreatedAt)
	if err != nil {
		return model.Orders{}, fmt.Errorf("insert order: %w", err)
	}
	const insertItem = `
		INSERT INTO order_items (
			order_id, product_id, quantity, unit_price_cents
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	for i := range order.Items {
		item := &order.Items[i]
		err := tx.QueryRowContext(ctx, insertItem, order.OrderID, item.ProductID, item.Quantity, item.UnitPriceCents).Scan(&item.ID)
		if err != nil {
			return model.Orders{}, fmt.Errorf("insert order item for product %d: %w", item.ProductID, err)
		}
		item.OrderID = order.OrderID
	}
	if err := tx.Commit(); err != nil {
		return model.Orders{}, fmt.Errorf("commit order transaction: %w", err)
	}
	return order, nil
}

func (r *PostgresOrderRepo) GetAll(ctx context.Context, limit, offset int) ([]model.Orders, error) {
	const selectOrdersQuery = `
		SELECT id, user_id, status, total_amount_cents, created_at
		FROM orders
		ORDER BY id desc
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, selectOrdersQuery, limit, offset)
	if err != nil {
		return []model.Orders{}, fmt.Errorf("select orders: %w", err)
	}
	defer rows.Close()
	var orders = make([]model.Orders, 0)
	orderIDs := make([]int, 0)
	orderIndex := make(map[int]int)

	for rows.Next() {
		var order model.Orders
		err = rows.Scan(&order.OrderID, &order.UserID, &order.Status, &order.TotalAmountCents, &order.CreatedAt)
		if err != nil {
			return []model.Orders{}, fmt.Errorf("scan order: %w", err)
		}

		order.Items = make([]model.OrderItem, 0)

		orderIndex[order.OrderID] = len(orders)
		orderIDs = append(orderIDs, order.OrderID)
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}
	if len(orderIDs) == 0 {
		return orders, nil
	}
	const selectItems = `
		SELECT id, order_id, product_id, quantity, unit_price_cents
		FROM order_items
		WHERE order_id = ANY($1)
		ORDER BY order_id, id 
	`
	itemrows, err := r.db.QueryContext(ctx, selectItems, pq.Array(orderIDs))
	if err != nil {
		return []model.Orders{}, fmt.Errorf("select items: %w", err)
	}
	defer itemrows.Close()
	for itemrows.Next() {
		var item model.OrderItem

		if err := itemrows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPriceCents,
		); err != nil {
			return []model.Orders{}, fmt.Errorf("scan order item: %w", err)
		}
		index, exists := orderIndex[item.OrderID]
		if !exists {
			return []model.Orders{}, fmt.Errorf("order item references unexpected order: %d", item.OrderID)
		}
		orders[index].Items = append(orders[index].Items, item)
	}

	if err := itemrows.Err(); err != nil {
		return []model.Orders{}, fmt.Errorf("iterate order items: %w", err)
	}
	return orders, nil
}

func (r *PostgresOrderRepo) GetAllByUserID(ctx context.Context, userId, limit, offset int) ([]model.Orders, error) {
	const queryOrders = `
		Select o.id, o.user_id, o.total_amount_cents, o.status, o.created_at
		FROM orders o
		WHERE o.user_id = $1
		ORDER BY o.id DESC
 		LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, queryOrders, userId, limit, offset)
	if err != nil {
		return []model.Orders{}, fmt.Errorf("query orders by userId: %w", err)
	}
	defer rows.Close()

	orders := make([]model.Orders, 0)
	orderIDs := make([]int, 0)
	orderIndex := make(map[int]int)

	for rows.Next() {
		var order model.Orders
		if err = rows.Scan(&order.OrderID, &order.UserID, &order.TotalAmountCents, &order.Status, &order.CreatedAt); err != nil {
			return []model.Orders{}, fmt.Errorf("scan order: %w", err)
		}

		order.Items = make([]model.OrderItem, 0)
		orderIndex[order.OrderID] = len(orders)
		orderIDs = append(orderIDs, order.OrderID)
		orders = append(orders, order)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}
	if len(orderIDs) == 0 {
		return orders, nil
	}
	const selectItems = `
		SELECT id, order_id, product_id, quantity, unit_price_cents
		FROM order_items
		WHERE order_id = ANY($1)
		ORDER BY order_id, id
`
	itemRows, err := r.db.QueryContext(ctx, selectItems, pq.Array(orderIDs))
	if err != nil {
		return nil, fmt.Errorf("select order items: %w", err)
	}

	defer itemRows.Close()

	for itemRows.Next() {
		var item model.OrderItem

		if err := itemRows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPriceCents,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		index, ok := orderIndex[item.OrderID]
		if !ok {
			return nil, fmt.Errorf("order item references unexpected order: %d", item.OrderID)
		}
		orders[index].Items = append(orders[index].Items, item)
	}
	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order: %w", err)
	}

	return orders, nil
}
func (r *PostgresOrderRepo) GetByID(ctx context.Context, id int) (model.Orders, error) {
	var order model.Orders
	err := r.db.QueryRowContext(ctx,
		`select o.id, o.user_id, o.total_amount_cents,  o.status, o.created_at from orders o WHERE o.id = $1`,
		id).Scan(&order.OrderID, &order.UserID, &order.TotalAmountCents, &order.Status, &order.CreatedAt)
	if err != nil {
		return model.Orders{}, err
	}
	rows, err := r.db.QueryContext(ctx, `select oi.id,  oi.product_id, oi.quantity, oi.unit_price_cents from order_items oi where oi.order_id = $1`, order.OrderID)
	if err != nil {
		return model.Orders{}, err
	}
	defer rows.Close()
	order.Items = []model.OrderItem{}
	for rows.Next() {
		var item model.OrderItem
		err = rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.UnitPriceCents)
		if err != nil {
			return model.Orders{}, err
		}
		item.OrderID = order.OrderID
		order.Items = append(order.Items, item)
	}
	return order, nil
}

// func (r *PostgresOrderRepo) Update(order model.Orders) (model.Orders, error) {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return model.Orders{}, err
// 	}
// 	defer tx.Rollback()
// 	_, err = tx.Exec(`Update orders set user_id = $1, status =$2 where id = $3 `, order.UserID, order.Status, order.ID)
// 	if err != nil {
// 		return model.Orders{}, err
// 	}
// 	_, err = tx.Exec(`DELETE from order_items where order_id = $1`, order.ID)
// 	if err != nil {
// 		return model.Orders{}, err
// 	}
// 	for i, items := range order.Items {
// 		var itemID int
// 		err = tx.QueryRow(`INSERT into order_items(order_id, product_id, quantity) VALUES ($1,$2,$3) RETURNING id`, order.ID, items.ProductID, items.Quantity).Scan(&itemID)
// 		if err != nil {
// 			return model.Orders{}, err
// 		}
// 		order.Items[i].ID = itemID
// 		order.Items[i].OrderId = order.ID
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return model.Orders{}, err
// 	}
// 	return order, nil
//

func (r *PostgresOrderRepo) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `delete from orders where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

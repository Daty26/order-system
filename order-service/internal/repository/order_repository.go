package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/Daty26/order-system/order-service/internal/db"
	"github.com/Daty26/order-system/order-service/internal/model"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
)

type OrderRep interface {
	Create(order model.Orders) (model.Orders, error)
	GetAll() ([]model.Orders, error)
	GetByID(id int) (model.Orders, error)
	Update(id int, productId int, quantity int, userId int) (model.Orders, error)
	Delete(id int) error
	GetAllByUserID(userId int) ([]model.Orders, error)
}
type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}
func (r *PostgresOrderRepo) Create(order model.Orders) (model.Orders, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return model.Orders{}, err
	}
	defer tx.Rollback()
	query := `insert into orders ( user_id, status) values ($1, $2) RETURNING id`
	err = tx.QueryRow(query, order.UserID, order.Status).Scan(&order.ID)
	if err != nil {
		return model.Orders{}, err
	}
	for i, item := range order.Items {
		var itemId int
		err = tx.QueryRow(`INSERT into order_items (order_id, product_id, quantity) VALUES ($1, $2,$3) RETURNING id`, order.ID, item.ProductID, item.Quantity).Scan(&itemId)
		if err != nil {
			return model.Orders{}, err
		}
		order.Items[i].ID = itemId
		order.Items[i].OrderId = order.ID
	}
	err = tx.Commit()
	if err != nil {
		return model.Orders{}, err
	}
	return order, nil
}
func (r *PostgresOrderRepo) GetAll() ([]model.Orders, error) {
	rows, err := r.db.Query(`Select o.id, o.user_id, o.status, o.created_at from orders o order by o.id desc`)
	if err != nil {
		return []model.Orders{}, err
	}
	defer rows.Close()
	var orders []model.Orders
	for rows.Next() {
		var order model.Orders
		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt)
		if err != nil {
			return []model.Orders{}, err
		}
		order.Items = []model.OrderItem{}
		itemRows, err := r.db.Query(`SELECT oi.id, oi.product_id, oi.quantity from order_items oi`)
		if err != nil {
			return []model.Orders{}, err
		}
		defer itemRows.Close()
		for itemRows.Next() {
			var item model.OrderItem
			err = itemRows.Scan(&item.ID, &item.ProductID, &item.Quantity)
			if err != nil {
				return []model.Orders{}, err
			}
			item.OrderId = order.ID
			order.Items = append(order.Items, item)
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (r *PostgresOrderRepo) GetAllByUserID(userId int) ([]model.Orders, error) {
	rows, err := r.db.Query(`Select o.id, o.user_id, o.status, o.created_at from orders o where o.user_id = $1 order by o.id desc`, userId)
	if err != nil {
		return []model.Orders{}, err
	}
	defer rows.Close()
	var orders []model.Orders
	for rows.Next() {
		var order model.Orders
		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt)
		if err != nil {
			return []model.Orders{}, err
		}
		order.Items = []model.OrderItem{}
		itemRows, err := r.db.Query(`SELECT oi.id, oi.product_id, oi.quantity from order_items oi where order_id = $1`, order.ID)
		if err != nil {
			return []model.Orders{}, err
		}
		for itemRows.Next() {
			var item model.OrderItem
			err = itemRows.Scan(&item.ID, &item.ProductID, &item.Quantity)
			if err != nil {
				return []model.Orders{}, err
			}
			item.OrderId = order.ID
			order.Items = append(order.Items, item)
		}
		itemRows.Close()
		orders = append(orders, order)
	}
	return orders, nil
}
func (r *PostgresOrderRepo) GetByID(id int) (model.Orders, error) {
	var order model.Orders
	err := r.db.QueryRow(`select o.id, o.user_id, o.status, o.created_at from orders o WHERE o.id = $1`, id).Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt)
	if err != nil {
		return model.Orders{}, err
	}
	rows, err := r.db.Query(`select oi.id  oi.product_id, oi.quantity from order_items oi where oi.order_id = $1`, order.ID)
	defer rows.Close()
	order.Items = []model.OrderItem{}
	for rows.Next() {
		var item model.OrderItem
		err = rows.Scan(&item.ID, &item.ProductID, &item.Quantity)
		if err != nil {
			return model.Orders{}, err
		}
		item.OrderId = order.ID
		order.Items = append(order.Items, item)
	}
	fmt.Println(order)
	return order, nil
}
func (r *PostgresOrderRepo) Update(id int, productId int, quantity int, userId int) (model.Orders, error) {
	var order model.Orders
	err := r.db.QueryRow(`update orders set product_id = $1, quantity = $2, user_id=$3 where id = $4 RETURNING id, product_id, quantity, user_id`, productId, quantity, userId, id).Scan(&order.ID, &order.ProductID, &order.Quantity, &order.UserID)
	if err != nil {
		return model.Orders{}, err
	}
	return order, err
}
func (r *PostgresOrderRepo) Delete(id int) error {
	res, err := r.db.Exec(`delete from orders where id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

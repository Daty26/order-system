package repository

import (
	"database/sql"
	"github.com/Daty26/order-system/order-service/internal/db"
	"github.com/Daty26/order-system/order-service/internal/model"
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
	rows, err := r.db.Query(`select o.id, o.user_id, o.status, oi.product_id, oi.order_id, oi.quantity, oi.id from orders o left join order_items oi on o.id = oi.order_id order by o.id desc`)
	if err != nil {
		return []model.Orders{}, err
	}
	defer rows.Close()

	var orders []model.Orders
	var currentOrder *model.Orders

	for rows.Next() {
		var (
			orderID   int
			userID    int
			status    string
			itemID    sql.NullInt64
			productID sql.NullInt64
			quantity  sql.NullInt64
		)

		err = rows.Scan(&orderID, &userID, &status, &itemID, &productID, &quantity)
		if err != nil {
			return nil, err
		}
		if currentOrder == nil || currentOrder.ID != orderID {
			order := model.Orders{
				ID:     orderID,
				UserID: userID,
				Status: status,
				Items:  []model.OrderItem{},
			}
			orders = append(orders, order)
			currentOrder = &orders[len(orders)-1]
		}
		if itemID.Valid {
			item := model.OrderItem{
				ID:        int(itemID.Int64),
				OrderId:   orderID,
				ProductID: int(productID.Int64),
				Quantity:  int(quantity.Int64),
			}
			currentOrder.Items = append(currentOrder.Items, item)
		}
	}
	return orders, nil
}
func (r *PostgresOrderRepo) GetAllByUserID(userId int) ([]model.Orders, error) {
	rows, err := r.db.Query(`select o.id, o.user_id, o.status, oi.product_id, oi.order_id, oi.quantity, oi.id from orders o left join order_items oi on o.id = oi.order_id where o.user_id = $1 order by o.id desc`, userId)
	if err != nil {
		return []model.Orders{}, err
	}
	defer rows.Close()

	var orders []model.Orders
	var currentOrder *model.Orders

	for rows.Next() {
		var (
			orderID   int
			userID    int
			status    string
			itemID    sql.NullInt64
			productID sql.NullInt64
			quantity  sql.NullInt64
		)

		err = rows.Scan(&orderID, &userID, &status, &itemID, &productID, &quantity)
		if err != nil {
			return nil, err
		}
		if currentOrder == nil || currentOrder.ID != orderID {
			order := model.Orders{
				ID:     orderID,
				UserID: userID,
				Status: status,
				Items:  []model.OrderItem{},
			}
			orders = append(orders, order)
			currentOrder = &orders[len(orders)-1]
		}
		if itemID.Valid {
			item := model.OrderItem{
				ID:        int(itemID.Int64),
				OrderId:   orderID,
				ProductID: int(productID.Int64),
				Quantity:  int(quantity.Int64),
			}
			currentOrder.Items = append(currentOrder.Items, item)
		}
	}
	return orders, nil
}
func (r *PostgresOrderRepo) GetByID(id int) (model.Orders, error) {
	var order model.Order
	err := r.db.QueryRow(`select id, product_id, quantity, user_id from orders WHERE id = $1`, id).Scan(&order.ID, &order.ProductID, &order.Quantity, &order.UserID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
func (r *PostgresOrderRepo) Update(id int, productId int, quantity int, userId int) (model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(`update orders set product_id = $1, quantity = $2, user_id=$3 where id = $4 RETURNING id, product_id, quantity, user_id`, productId, quantity, userId, id).Scan(&order.ID, &order.ProductID, &order.Quantity, &order.UserID)
	if err != nil {
		return model.Order{}, err
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

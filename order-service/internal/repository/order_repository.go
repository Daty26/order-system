package repository

import (
	"database/sql"
	"github.com/Daty26/order-system/order-service/internal/model"
)

type OrderRep interface {
	Create(order model.Order) (model.Order, error)
	GetAll() ([]model.Order, error)
	GetByID(id int) (model.Order, error)
	Update(id int, item string, quantity int, userId int) (model.Order, error)
	Delete(id int) error
	GetAllByUserID(userId int) ([]model.Order, error)
}
type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}
func (r *PostgresOrderRepo) Create(order model.Order) (model.Order, error) {
	query := `insert into orders (item, quantity, user_id) values ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, order.Item, order.Quantity, order.UserID).Scan(&order.ID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
func (r *PostgresOrderRepo) GetAll() ([]model.Order, error) {
	rows, err := r.db.Query(`select id, item, quantity from orders`)
	if err != nil {
		return []model.Order{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.Item, &o.Quantity); err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
func (r *PostgresOrderRepo) GetAllByUserID(userId int) ([]model.Order, error) {
	rows, err := r.db.Query(`select id, item, quantity, user_id from orders where user_id = $1`, userId)
	if err != nil {
		return []model.Order{}, err
	}
	defer rows.Close()
	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.ID, &order.Item, &order.Quantity, &order.UserID)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (r *PostgresOrderRepo) GetByID(id int) (model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(`select id, item, quantity, user_id from orders WHERE id = $1`, id).Scan(&order.ID, &order.Item, &order.Quantity, &order.UserID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
func (r *PostgresOrderRepo) Update(id int, item string, quantity int, userId int) (model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(`update orders set item = $1, quantity = $2, user_id=$3 where id = $4 RETURNING id, item, quantity, user_id`, item, quantity, userId, id).Scan(&order.ID, &order.Item, &order.Quantity, &order.UserID)
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

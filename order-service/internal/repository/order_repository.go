package repository

import (
	"database/sql"
	"github.com/Daty26/order-system/order-service/internal/db"
	"github.com/Daty26/order-system/order-service/internal/model"
)

type OrderRep interface {
	Create(order model.Order) (model.Order, error)
	GetAll() ([]model.Order, error)
	GetByID(id int) (model.Order, error)
	Update(id int, item string, amount int) (model.Order, error)
	Delete(id int) error
}
type PostgresOrderRepo struct{}

func NewRepo() *PostgresOrderRepo {
	return &PostgresOrderRepo{}
}
func (r *PostgresOrderRepo) Create(order model.Order) (model.Order, error) {
	query := `insert into orders (item, amount) values ($1, $2) RETURNING id`
	err := db.DataB.QueryRow(query, order.Item, order.Amount).Scan(&order.ID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
func (r *PostgresOrderRepo) GetAll() ([]model.Order, error) {
	rows, err := db.DataB.Query(`select id, item, amount from orders`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.Item, &o.Amount); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
func (r *PostgresOrderRepo) GetByID(id int) (model.Order, error) {
	var order model.Order
	err := db.DataB.QueryRow(`select id, item, amount from orders WHERE id = $1`, id).Scan(&order.ID, &order.Item, &order.Amount)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
func (r *PostgresOrderRepo) Update(id int, item string, amount int) (model.Order, error) {
	var order model.Order
	err := db.DataB.QueryRow(`update orders set item = $1, amount = $2 where id = $3 RETURNING id, item, amount`, item, amount, id).Scan(&order.ID, &order.Item, &order.Amount)
	if err != nil {
		return model.Order{}, err
	}
	return order, err
}
func (r *PostgresOrderRepo) Delete(id int) error {
	res, err := db.DataB.Exec(`delete from orders where id = $1`, id)
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

package repository

import (
	"database/sql"

	"github.com/Daty26/order-system/payment-service/internal/model"
)

type PaymentRep interface {
	Save(payment model.Payment) (model.Payment, error)
	GetAll() ([]model.Payment, error)
}
type PostgresPaymentRep struct {
	db *sql.DB
}

func NewPostgresRep(db *sql.DB) *PostgresPaymentRep {
	return &PostgresPaymentRep{db: db}
}

func (r *PostgresPaymentRep) Save(payment model.Payment) (model.Payment, error) {
	query := `INSERT INTO payments(order_id, status, amount) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, payment.OrderID, payment.Status, payment.Amount).Scan(&payment.ID)
	if err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresPaymentRep) GetAll() ([]model.Payment, error) {
	rows, err := r.db.Query(`SELECT id, order_id, status, amount from payments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		err := rows.Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.Amount)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *PostgresPaymentRep) GetByID(id int) (model.Payment, error) {
	var payment model.Payment
	err := r.db.QueryRow("SELECT id, order_id, status, amount from payment where id=$1", id).Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.Amount)
	if err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daty26/order-system/payment-service/internal/model"
)

type PaymentRep interface {
	Save(ctx context.Context, params ProcessPaymentParams) (model.Payment, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Payment, error)
	GetByID(ctx context.Context, id int) (model.Payment, error)
	Update(ctx context.Context, id int, status model.PaymentStatus, amount float64) (model.Payment, error)
	Delete(ctx context.Context, id int) error
	GetAllByUserId(ctx context.Context, userId int) ([]model.Payment, error)
}

type PostgresPaymentRep struct {
	db *sql.DB
}

func NewPostgresRep(db *sql.DB) *PostgresPaymentRep {
	return &PostgresPaymentRep{db: db}
}

func (r *PostgresPaymentRep) Save(ctx context.Context, params ProcessPaymentParams) (model.Payment, error) {
	query := `
		INSERT INTO payments(order_id, status, amount_cents, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, order_id, status, amount_cents, user_id`
	var payment model.Payment
	if err := r.db.QueryRowContext(
		ctx,
		query,
		params.OrderID,
		params.Status,
		params.AmountCents,
		params.UserID,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Status,
		&payment.AmountCents,
		&payment.UserID,
	); err != nil {
		return model.Payment{}, fmt.Errorf("query insert payment: %w", err)
	}
	return payment, nil
}

func (r *PostgresPaymentRep) GetAll(ctx context.Context, limit, offset int) ([]model.Payment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, order_id, status, amount, user_id from payments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		err := rows.Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.AmountCents, &payment.UserID)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *PostgresPaymentRep) GetAllByUserId(ctx context.Context, userId int) ([]model.Payment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, order_id, status, amount, user_id from payments where user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		err := rows.Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.AmountCents, &payment.UserID)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *PostgresPaymentRep) GetByID(ctx context.Context, id int) (model.Payment, error) {
	var payment model.Payment
	err := r.db.QueryRowContext(ctx, "SELECT id, order_id, status, amount, user_id from payments where id=$1", id).Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.AmountCents, &payment.UserID)
	if err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresPaymentRep) Update(ctx context.Context, id int, status model.PaymentStatus, amount float64) (model.Payment, error) {
	var payment model.Payment
	query := `update payments SET status=$1, amount=$2 where id = $3 RETURNING id,order_id, status, amount, user_id`
	if err := r.db.QueryRowContext(ctx, query, status, amount, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Status,
		&payment.AmountCents,
		&payment.UserID,
	); err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresPaymentRep) Delete(ctx context.Context, id int) error {
	query := `Delete from payments where id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
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

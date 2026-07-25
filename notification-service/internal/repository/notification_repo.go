package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daty26/order-system/notification-service/internal/model"
)

type NotificationRepo interface {
	Insert(ctx context.Context, params InsertParams) (model.Notification, error)
	GetAllByUserID(ctx context.Context, params GetAllByUserIDParams) ([]model.Notification, error)
	GetByID(ctx context.Context, id int) (model.Notification, error)
	GetByStatus(ctx context.Context, params GetByStatusParams) ([]model.Notification, error)
	UpdateStatusByID(ctx context.Context, params UpdateStatusByIDParams) (model.Notification, error)
	DeleteByID(ctx context.Context, id int) error
	GetAll(ctx context.Context, limit, offset int) ([]model.Notification, error)
}

type PostgresNotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{db: db}
}

func (r *PostgresNotificationRepo) Insert(ctx context.Context, params InsertParams) (model.Notification, error) {
	query := `
		INSERT INTO notifications (order_id, payment_id, status, message, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, order_id, payment_id, status, message, user_id, created_at`
	var notification model.Notification
	row := r.db.QueryRowContext(
		ctx,
		query,
		params.OrderID,
		params.PaymentID,
		params.Status,
		params.Message,
		params.UserID,
	)
	err := row.Scan(
		&notification.ID,
		&notification.OrderID,
		&notification.PaymentID,
		&notification.Status,
		&notification.Message,
		&notification.UserID,
		&notification.CreatedAt,
	)
	if err != nil {
		return model.Notification{}, fmt.Errorf("insert notification: %w", err)
	}
	return notification, nil
}

func (r *PostgresNotificationRepo) GetAllByUserID(ctx context.Context, params GetAllByUserIDParams) ([]model.Notification, error) {
	query := `
		SELECT id, order_id, payment_id, status, message, user_id, created_at
		FROM notifications
		where user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, query, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return []model.Notification{}, fmt.Errorf("query notification: %w", err)
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		if err = rows.Scan(&notification.ID, &notification.OrderID, &notification.PaymentID, &notification.Status, &notification.Message, &notification.UserID, &notification.CreatedAt); err != nil {
			return []model.Notification{}, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return []model.Notification{}, fmt.Errorf("rows iterate: %w", err)
	}
	return notifications, nil
}

func (r *PostgresNotificationRepo) GetAll(ctx context.Context, limit, offset int) ([]model.Notification, error) {
	query := `
		SELECT id, order_id, payment_id, status, message, user_id, created_at
		FROM notifications
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return []model.Notification{}, fmt.Errorf("query notification: %w", err)
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		if err = rows.Scan(
			&notification.ID,
			&notification.OrderID,
			&notification.PaymentID,
			&notification.Status,
			&notification.Message,
			&notification.UserID,
			&notification.CreatedAt,
		); err != nil {
			return []model.Notification{}, fmt.Errorf("scan rows: %w", err)
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return []model.Notification{}, fmt.Errorf("iterate rows: %w", err)
	}
	return notifications, nil
}

func (r *PostgresNotificationRepo) GetByID(ctx context.Context, id int) (model.Notification, error) {
	var notification model.Notification
	query := `
		SELECT id, order_id, payment_id, status, message, user_id, created_at
		from notifications
		where id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.OrderID,
		&notification.PaymentID,
		&notification.Status,
		&notification.Message,
		&notification.UserID,
		&notification.CreatedAt,
	)
	if err != nil {
		return model.Notification{}, fmt.Errorf("select notification by id: %w", err)
	}
	return notification, nil
}

func (r *PostgresNotificationRepo) GetByStatus(ctx context.Context, params GetByStatusParams) ([]model.Notification, error) {
	query := `
		SELECT id, order_id, payment_id, status, message, user_id, created_at
		from notifications
		where status=$1 and user_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
`
	rows, err := r.db.QueryContext(ctx, query, params.Status, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return []model.Notification{}, fmt.Errorf("get notifications by status: %w", err)
	}
	defer rows.Close()
	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		if err = rows.Scan(
			&notification.ID,
			&notification.OrderID,
			&notification.PaymentID,
			&notification.Status,
			&notification.Message,
			&notification.UserID,
			&notification.CreatedAt,
		); err != nil {
			return []model.Notification{}, fmt.Errorf("scan notifications: %w", err)
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return []model.Notification{}, fmt.Errorf("iterate rows: %w", err)
	}
	return notifications, nil
}

func (r *PostgresNotificationRepo) UpdateStatusByID(ctx context.Context, params UpdateStatusByIDParams) (model.Notification, error) {
	var notification model.Notification
	query := `
		UPDATE notifications
		SET status = $1
		WHERE id = $2
		RETURNING id, order_id, payment_id, status, message, user_id, created_at`
	err := r.db.QueryRowContext(ctx, query, params.Status, params.ID).Scan(&notification.ID, &notification.OrderID, &notification.PaymentID, &notification.Status, &notification.Message, &notification.UserID, &notification.CreatedAt)
	if err != nil {
		return model.Notification{}, fmt.Errorf("update notification by id: %w", err)
	}
	return notification, nil
}

func (r *PostgresNotificationRepo) DeleteByID(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `Delete from notifications where id =$1`, id)
	if err != nil {
		return fmt.Errorf("delete notification: %w", err)
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

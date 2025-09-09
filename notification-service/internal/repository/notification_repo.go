package repository

import (
	"database/sql"

	"github.com/Daty26/order-system/notification-service/internal/model"
)

type NotificationRepo interface {
	Insert(notification model.Notification) (model.Notification, error)
}
type PostgresNotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{db: db}
}

func (nf *PostgresNotificationRepo) Insert(notification model.Notification) (model.Notification, error) {
	query := `INSERT INTO notifications (order_id, payment_id, status, message) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	row := nf.db.QueryRow(query, notification.OrderID, notification.PaymentID, notification.Status, notification.Message)
	err := row.Scan(&notification.ID, &notification.CreatedAt)
	if err != nil {
		return model.Notification{}, err
	}
	return notification, nil
}

package repository

import (
	"database/sql"

	"github.com/Daty26/order-system/notification-service/internal/model"
)

type NotificationRepo interface {
	Insert(notification model.Notification) (model.Notification, error)
	GetAllByUserID(userId int) ([]model.Notification, error)
	GetByID(id int) (model.Notification, error)
	GetByStatus(status model.NotificationStatus, userId int) ([]model.Notification, error)
	UpdateStatusByID(id int, status model.NotificationStatus) (model.Notification, error)
	DeleteByID(id int) error
}
type PostgresNotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{db: db}
}

func (nf *PostgresNotificationRepo) Insert(notification model.Notification) (model.Notification, error) {
	query := `INSERT INTO notifications (order_id, payment_id, status, message, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	row := nf.db.QueryRow(query, notification.OrderID, notification.PaymentID, notification.Status, notification.Message, notification.UserID)

	err := row.Scan(&notification.ID, &notification.CreatedAt)
	if err != nil {
		return model.Notification{}, err
	}
	return notification, nil
}
func (nf *PostgresNotificationRepo) GetAllByUserID(userId int) ([]model.Notification, error) {
	rows, err := nf.db.Query(`SELECT id, order_id, payment_id, status, message, user_id, created_at FROM notifications where user_id = $1 ORDER BY created_at DESC`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		if err = rows.Scan(&notification.ID, &notification.OrderID, &notification.PaymentID, &notification.Status, &notification.Message, &notification.UserID, &notification.CreatedAt); err != nil {
			return notifications, err
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return notifications, err
	}
	return notifications, nil
}
func (nf *PostgresNotificationRepo) GetAl() ([]model.Notification, error) {
	rows, err := nf.db.Query(`SELECT id, order_id, payment_id, status, message, user_id, created_at FROM notifications ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		if err = rows.Scan(&notification.ID, &notification.OrderID, &notification.PaymentID, &notification.Status, &notification.Message, &notification.UserID, &notification.CreatedAt); err != nil {
			return notifications, err
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return notifications, err
	}
	return notifications, nil
}
func (nf *PostgresNotificationRepo) GetByID(id int) (model.Notification, error) {
	var notification model.Notification
	query := `SELECT id, order_id, payment_id, status, message, user_id,created_at from notifications where id = $1`
	err := nf.db.QueryRow(query, id).Scan(&notification.ID, &notification.OrderID, &notification.PaymentID, &notification.Status, &notification.Message, &notification.UserID, &notification.CreatedAt)
	if err != nil {
		return model.Notification{}, err
	}
	return notification, nil
}
func (nf *PostgresNotificationRepo) GetByStatus(status model.NotificationStatus, userId int) ([]model.Notification, error) {
	rows, err := nf.db.Query(`SELECT id, order_id, payment_id, status, message, user_id, created_at from notifications where status=$1 and user_id = $2`, status, userId)
	if err != nil {
		return []model.Notification{}, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		if err = rows.Scan(&notification.ID, &notification.OrderID, &notification.PaymentID, &notification.Status, &notification.Message, &notification.UserID, &notification.CreatedAt); err != nil {
			return notifications, err
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return notifications, err
	}
	return notifications, err
}
func (nf *PostgresNotificationRepo) UpdateStatusByID(id int, status model.NotificationStatus) (model.Notification, error) {
	var notification model.Notification
	query := `update notifications  SET status = $1 where id = $2 RETURNING id, order_id, payment_id, status, message,user_id, created_at`
	err := nf.db.QueryRow(query, status, id).Scan(&notification.ID, &notification.OrderID, &notification.PaymentID, &notification.Status, &notification.Message, &notification.UserID, &notification.CreatedAt)
	if err != nil {
		return model.Notification{}, err
	}
	return notification, err
}
func (nf *PostgresNotificationRepo) DeleteByID(id int) error {
	res, err := nf.db.Exec(`Delete from notifications where id =$1`, id)
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

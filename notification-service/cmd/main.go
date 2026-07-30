package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Daty26/order-system/notification-service/internal/middleware"

	"github.com/Daty26/order-system/notification-service/internal/api"
	"github.com/Daty26/order-system/notification-service/internal/db"
	"github.com/Daty26/order-system/notification-service/internal/kafka"
	"github.com/Daty26/order-system/notification-service/internal/model"
	"github.com/Daty26/order-system/notification-service/internal/repository"
	"github.com/Daty26/order-system/notification-service/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	db.InitDB()
	defer db.DataDB.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	rep := repository.NewNotificationRepo(db.DataDB)
	serv := service.NewNotificationService(rep)
	handler := api.NewNotificationHardler(serv, logger)

	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")

	consumer, err := kafka.NewKafkaConsumer(kafkaBrokers, func(value []byte) {
		handlePaymentCompleted(value, serv, logger)
	})
	if err != nil {
		logger.Error("failed to create kafka consumer", "error", err)
		os.Exit(1)
	}

	defer consumer.Close()
	if err := consumer.Consume("payment.completed"); err != nil {
		logger.Error("failed to consume payment.completed", "error", err)
		os.Exit(1)
	}
	r := chi.NewRouter()
	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("notification-service is working"))
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/notifications", handler.InsertNotification)
		r.Get("/notifications", handler.GetNotifications)
		r.Get("/notifications/{id}", handler.GetNotificationByID)
		r.Get("/notifications/status/{status}", handler.GetNotificationsByStatus)
		r.Put("/notifications/{id}/status", handler.UpdateNotificationStatus)
		r.Delete("/notifications/{id}", handler.DeleteNotificationByID)
	})

	if err = http.ListenAndServe(":8083", r); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func handlePaymentCompleted(
	value []byte,
	serv *service.NotificationService,
	logger *slog.Logger,
) {
	var event model.PaymentCreated
	if err := json.Unmarshal(value, &event); err != nil {
		logger.Warn("invalid payment.completed event", "error", err)
		return
	}
	input := service.InsertInput{
		OrderID:   event.OrderID,
		PaymentID: event.PaymentID,
		Status:    model.NotificationSent,
		UserID:    event.UserID,
		Message:   "payment has been completed",
	}
	notification, err := serv.Insert(context.Background(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotificationAlreadyExists):
			logger.Error(
				"notification already exists",
				"payment_id", event.PaymentID,
				"order_id", event.OrderID,
			)
		default:
			logger.Error(
				"failed to create notification from payment.completed",
				"error", err,
				"order_id", event.OrderID,
				"payment_id", event.PaymentID,
				"user_id", event.UserID,
			)
		}
		return
	}
	logger.Info(
		"notification created",
		"notification_id", notification.ID,
		"order_id", notification.OrderID,
		"payment_id", notification.PaymentID,
	)
}

package main

import (
	"encoding/json"
	"github.com/Daty26/order-system/auth/middleware"
	"log"
	"net/http"

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
	rep := repository.NewNotificationRepo(db.DataDB)
	serv := service.NewNotificationService(rep)
	handler := api.NewNotificationHandler(serv)

	consumePaymentCreated := func(value []byte) {
		log.Println(string(value))
		var paymentCreated model.PaymentCreated
		if err := json.Unmarshal(value, &paymentCreated); err != nil {
			log.Println("Kafka handler, notification consumer: " + err.Error())
			return
		}
		log.Println(paymentCreated)
		if _, err := serv.Insert(paymentCreated.OrderID, paymentCreated.PaymentID, model.NotificationSent, paymentCreated.UserID, "payment has been created"); err != nil {
			log.Println("Kafka process for notification service failed: ", err)
			return
		}
		log.Printf("Notification is created for orderid=%v, paymentid=%v", paymentCreated.OrderID, paymentCreated.PaymentID)
	}
	consumer, err := kafka.NewKafkaConsumer([]string{"localhost:9092"}, consumePaymentCreated)
	if err != nil {
		log.Printf("Failed to create new kafka consumer: %v", err.Error())
		return
	}
	go func() {
		log.Println("Consuming topic payment.completed")
		if err := consumer.Consume("payment.completed"); err != nil {
			log.Printf("kafka consumer error: %v", err)
		}
	}()
	defer consumer.Close()

	r := chi.NewRouter()
	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("notification-service is working"))
		if err != nil {
			return
		}
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/notifications", handler.InsertNotification)
		r.Get("/notifications", handler.GetNotifications)
		r.Get("/notifications/{id}", handler.GetNotificationByID)
		r.Get("/notifications/status/{status}", handler.GetNotificationsByStatus)
		r.Put("/notifications/{id}/status", handler.UpdateNotificationStatusByID)
		r.Delete("/notifications/{id}", handler.DeleteNotificationByID)
	})
	err = http.ListenAndServe(":8083", r)
	if err != nil {
		log.Fatal(err.Error())
	}
}

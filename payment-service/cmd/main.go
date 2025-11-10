package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Daty26/order-system/payment-service/internal/kafka"

	_ "github.com/Daty26/order-system/payment-service/docs"
	"github.com/Daty26/order-system/payment-service/internal/api"
	"github.com/Daty26/order-system/payment-service/internal/db"
	"github.com/Daty26/order-system/payment-service/internal/repository"
	"github.com/Daty26/order-system/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @description Payment service for the order system
// @host localhost:8080
// @BasePath /
func main() {
	db.InitDB()
	defer db.DataDB.Close()

	producer, err := kafka.NewKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("couldn't create producer: %v", err)
	}
	defer producer.Close()

	repo := repository.NewPostgresRep(db.DataDB)
	srv := service.NewPaymentService(repo, producer)

	type orderCreated struct {
		OrderID int `json:"order_id"`
		Items   interface {
		} `json:"items"`
		Quantity int `json:"quantity"`
	}

	consumeOrderCreated := func(value []byte) {
		var order orderCreated
		if err := json.Unmarshal(value, &order); err != nil {
			log.Println("kafka handler: bad payload:", err)
			return
		}
		if _, err := srv.ProcessPayment(order.PaymentId, order.Amount); err != nil {
			log.Println("kafka handler: process payment failed:", err)
			return
		}
		log.Printf("Processed payment for order %d\n", order.PaymentId)
	}
	consumer, err := kafka.NewKafkaConsumer([]string{"localhost:9092"}, consumeOrderCreated)
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	if err = consumer.Consume("order.created"); err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}
	defer consumer.Close()
	handler := api.NewRepoPyament(srv)

	r := chi.NewRouter()

	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("payment-service is working"))
		if err != nil {
			return
		}
	})
	r.Get("/payments", handler.GetPayments)
	r.Put("/payments/{id}", handler.UpdatePayment)
	r.Get("/payments/{id}", handler.GetPaymentByID)
	r.Post("/payments", handler.CreatePayment)
	r.Delete("/payments/{id}", handler.DeletePayment)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	fmt.Println("starting payment service on port 8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

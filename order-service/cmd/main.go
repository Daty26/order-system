package main

import (
	"github.com/Daty26/order-system/auth/middleware"
	_ "github.com/Daty26/order-system/order-service/docs"
	"github.com/Daty26/order-system/order-service/internal/api"
	"github.com/Daty26/order-system/order-service/internal/db"
	"github.com/Daty26/order-system/order-service/internal/kafka"
	"github.com/Daty26/order-system/order-service/internal/repository"
	"github.com/Daty26/order-system/order-service/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// @title Order Service API
// @version 1.0
// @description This is the order service for the Event-Driven Order System.
// @host localhost:8080
// @BasePath /
func main() {
	db.InitDB()
	defer db.DataB.Close()

	prod, err := kafka.NewKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalln("failed to create Kafka producer: " + err.Error())
	}
	defer prod.Close()

	repo := repository.NewPostgresRepo(db.DataB)

	svc := service.NewOrderService(repo, prod)

	handler := api.NewOrderHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("order-service is healthy"))
		if err != nil {
			return
		}
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Put("/orders/{id}", handler.UpdateOrder)
		r.Delete("/orders/{id}", handler.DeleteOrder)
		r.Post("/orders", handler.CreateOrder)
		r.Get("/orders", handler.GetOrders)
		r.Get("/orders/{id}", handler.GetOrderByID)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("starting order-system on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("couldn't start the server: " + err.Error())
		return
	}
}

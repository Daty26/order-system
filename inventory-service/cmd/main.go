package main

import (
	"github.com/Daty26/order-system/auth/middleware"
	"github.com/Daty26/order-system/inventory-service/internal/api"
	"github.com/Daty26/order-system/inventory-service/internal/db"
	"github.com/Daty26/order-system/inventory-service/internal/kafka"
	"github.com/Daty26/order-system/inventory-service/internal/repository"
	"github.com/Daty26/order-system/inventory-service/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	db.InitDB()
	defer db.DataDB.Close()
	repo := repository.NewPostgresInventoryRepo(db.DataDB)
	svc := service.NewInventoryService(repo)
	handler := api.NewInventoryHandler(svc)
	r := chi.NewRouter()
	go func() {
		consumer, err := kafka.NewKafkaConsumer([]string{"localhost:9092"}, svc)
		if err != nil {
			log.Fatalf("couldn't start consumer: %s" + err.Error())
		}
		if err = consumer.Consume("order.created"); err != nil {
			log.Fatalf("couldn't consume the topic: " + err.Error())
		}
	}()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("inventory-service is healthy"))
		if err != nil {
			return
		}
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/products", handler.InsertProduct)
		r.Put("/products/{id}", handler.UpdateQuantity)
	})
	r.Get("/products", handler.GetAllProducts)
	err := http.ListenAndServe(":8084", r)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
}

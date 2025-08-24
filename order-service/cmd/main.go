package main

import (
	"database/sql"
	"fmt"
	_ "github.com/Daty26/order-system/order-service/docs"
	"github.com/Daty26/order-system/order-service/internal/api"
	"github.com/Daty26/order-system/order-service/internal/db"
	"github.com/Daty26/order-system/order-service/internal/repository"
	"github.com/Daty26/order-system/order-service/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
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
	defer func(DataB *sql.DB) {
		err := DataB.Close()
		if err != nil {

		}
	}(db.DataB)

	repo := repository.NewRepo()

	svc := service.NewOrderService(repo)

	handler := api.NewOrderHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("order-service is healthy"))
		if err != nil {
			return
		}
	})

	r.Get("/orders", handler.GetOrders)
	r.Get("/orders/{id}", handler.GetOrderByID)
	r.Put("/orders/{id}", handler.UpdateOrder)
	r.Delete("/orders/{id}", handler.DeleteOrder)
	r.Post("/orders", handler.CreateOrder)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	fmt.Println("starting order-system on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}

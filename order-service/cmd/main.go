package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/Daty26/order-system/order-service/docs"
	"github.com/Daty26/order-system/order-service/internal/api"
	"github.com/Daty26/order-system/order-service/internal/client/inventory"
	"github.com/Daty26/order-system/order-service/internal/db"
	"github.com/Daty26/order-system/order-service/internal/kafka"
	"github.com/Daty26/order-system/order-service/internal/middleware"
	"github.com/Daty26/order-system/order-service/internal/repository"
	"github.com/Daty26/order-system/order-service/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	prod, err := kafka.NewKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatalln("failed to create Kafka producer: " + err.Error())
	}

	defer prod.Close()

	inventoryClient := inventory.NewClient(
		getEnv("INVENTORY_SERVICE_URL", "http://localhost:8084"),
		&http.Client{Timeout: 2 * time.Second},
	)

	repo := repository.NewPostgresRepo(db.DataB)

	svc := service.NewOrderService(repo, prod, inventoryClient)

	handler := api.NewOrderHandler(svc, logger)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("order-service is healthy"))
		if err != nil {
			return
		}
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		// r.Put("/orders/{id}", handler.UpdateOrder)
		r.Delete("/orders/{id}", handler.DeleteOrder)
		r.Post("/orders", handler.CreateOrder)
		r.Get("/orders", handler.GetOrders)
		r.Get("/orders/{id}", handler.GetOrderByID)

		// change order status
		r.Patch("/orders/{id}/cancel", handler)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("starting order-system on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("couldn't start the server: " + err.Error())
		return
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Daty26/order-system/payment-service/internal/client/order"
	"github.com/Daty26/order-system/payment-service/internal/middleware"

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

	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")

	producer, err := kafka.NewKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatalf(" create producer: %v", err)
	}
	defer producer.Close()

	orderClient := order.NewClient(
		getEnv("ORDER_SERVICE_URL", "http://localhost:8080"),
		&http.Client{Timeout: 2 * time.Second},
	)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	repo := repository.NewPostgresRep(db.DataDB)
	srv := service.NewPaymentService(repo, producer, orderClient)
	handler := api.NewPaymentHandler(srv, logger)

	r := chi.NewRouter()

	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("payment-service is working"))
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Put("/payments/{id}", handler.UpdatePayment)
		r.Get("/payments/{id}", handler.GetPaymentByID)
		r.Post("/payments", handler.CreatePayment)
		r.Delete("/payments/{id}", handler.DeletePayment)
		r.Get("/payments", handler.GetPayments)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("starting payment service on port 8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

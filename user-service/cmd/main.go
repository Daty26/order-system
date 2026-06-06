package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Daty26/order-system/user-service/internal/db"
	"github.com/Daty26/order-system/user-service/internal/middleware"
	"github.com/Daty26/order-system/user-service/internal/repository"
	"github.com/Daty26/order-system/user-service/internal/service"
	transport_http_handler "github.com/Daty26/order-system/user-service/internal/transport/http/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	db.DBInit()
	defer db.DBConn.Close()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	repo := repository.NewPostgresRepository(db.DBConn)
	srv := service.NewUserService(repo, jwtSecret)
	handler := transport_http_handler.NewUserHandler(srv, logger)

	r := chi.NewRouter()
	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("user-service is working"))
		if err != nil {
			return
		}
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/user/me", handler.Me)
	})
	r.Post("/user/register", handler.CreateUser)
	r.Post("/user/login", handler.LoginUser)
	r.Get("/users", handler.GetAll)
	log.Println("starting user service on port 8085")
	err := http.ListenAndServe(":8085", r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
		return
	}
}

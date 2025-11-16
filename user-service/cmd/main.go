package main

import (
	"github.com/Daty26/order-system/auth/middleware"
	"github.com/Daty26/order-system/user-service/internal/api"
	"github.com/Daty26/order-system/user-service/internal/db"
	"github.com/Daty26/order-system/user-service/internal/repository"
	"github.com/Daty26/order-system/user-service/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	db.DBInit()
	defer db.DBConn.Close()
	repo := repository.NewPostgresRepository(db.DBConn)
	srv := service.NewUserService(repo)
	handler := api.NewUserHandler(srv)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/user/me", handler.Me)
	})
	r.Post("/user/register", handler.CreateUser)
	r.Post("/user/login", handler.Login)
	log.Println("starting user service on port 8085")
	err := http.ListenAndServe(":8085", r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
		return
	}

}

package main

import (
	"fmt"
	"github.com/Daty26/order-system/notification-service/internal/api"
	"github.com/Daty26/order-system/notification-service/internal/db"
	"github.com/Daty26/order-system/notification-service/internal/repository"
	"github.com/Daty26/order-system/notification-service/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	db.InitDB()
	defer db.DataDB.Close()
	fmt.Println(1)
	rep := repository.NewNotificationRepo(db.DataDB)
	fmt.Println(2)
	serv := service.NewNotificationService(rep)
	handler := api.NewNotificationHandler(serv)
	r := chi.NewRouter()
	r.Post("/notification", handler.InsertNotification)
	err := http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}

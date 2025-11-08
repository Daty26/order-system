package main

import (
	"github.com/Daty26/order-system/inventory-service/internal/api"
	"github.com/Daty26/order-system/inventory-service/internal/db"
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
	r.Get("/products", handler.GetAllProducts)
	r.Post("/products", handler.InsertProduct)
	r.Put("/products/{id}", handler.UpdateQuantity)
	err := http.ListenAndServe(":8084", r)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
}

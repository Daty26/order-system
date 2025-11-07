package api

import (
	"encoding/json"
	"github.com/Daty26/order-system/inventory-service/internal/model"
	"github.com/Daty26/order-system/inventory-service/internal/service"
	"net/http"
)

type InventoryHandler struct {
	serv *service.InventoryService
}

func NewInventoryHandler(serv *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{serv: serv}
}
func (ih *InventoryHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := ih.serv.GetAll()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(w, http.StatusOK, products)
}
func (ih *InventoryHandler) InsertProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Incorrect req format: "+err.Error())
		return
	}
	product := model.Product{
		Name:     req.Name,
		Quantity: req.Quantity,
	}
	productCrteated, err := ih.serv.Insert(product)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(w, http.StatusCreated, productCrteated)

}

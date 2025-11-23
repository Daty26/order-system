package api

import (
	"encoding/json"
	"github.com/Daty26/order-system/inventory-service/internal/model"
	"github.com/Daty26/order-system/inventory-service/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
	if r.Context().Value("role") != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to add products")
	}
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
func (ih *InventoryHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Couldn't convert string to int: "+err.Error())
		return
	}
	if r.Context().Value("role") != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to change product quantity")
		return
	}
	var req struct {
		Quantity int `json:"quantity"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	quantity, err := ih.serv.UpdateQuantity(id, req.Quantity)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(w, http.StatusOK, quantity)
}

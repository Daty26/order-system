package api

import (
	"encoding/json"
	"github.com/Daty26/order-system/inventory-service/internal/model"
	"github.com/Daty26/order-system/inventory-service/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
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
	limit, offset, ok := parsePagination(r)
	if !ok {
		ErrorResponse(w, http.StatusBadRequest, "invalid pagination params")
		return
	}
	products, err := ih.serv.GetAll(r.Context(), limit, offset)
	if err != nil {
		log.Println(err.Error())
		ErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	SuccessResponse(w, http.StatusOK, products)
}

func (ih *InventoryHandler) InsertProduct(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to add products")
		return
	}
	var req InsertProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		ErrorResponse(w, http.StatusBadRequest, "Incorrect req body")
		return
	}
	insertInput := service.InsertProductInput{
		Name:       req.Name,
		Quantity:   req.Quantity,
		PriceCents: req.PriceCents,
	}
	productCreated, err := ih.serv.InsertProduct(r.Context(), insertInput)
	if err != nil {
		log.Printf("failed to insert product: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	SuccessResponse(w, http.StatusCreated, productCreated)
}

func (ih *InventoryHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
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
		ErrorResponse(w, http.StatusBadRequest, "invalid req body")
		return
	}
	quantity, err := ih.serv.UpdateQuantity(id, req.Quantity)
	if err != nil {
		log.Printf("Couldn't update quantity: %s", err.Error())
		ErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	SuccessResponse(w, http.StatusOK, quantity)
}

func (ih *InventoryHandler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Context().Value("role") != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to change product price")
		return
	}
	var req struct {
		Price float64 `json:"price"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid req body")
		return
	}
	price, err := ih.serv.UpdatePrice(id, req.Price)
	if err != nil {
		log.Printf("Couldn't update price: %s", err.Error())
		ErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	SuccessResponse(w, http.StatusOK, price)
}

package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Daty26/order-system/inventory-service/internal/service"
	"github.com/go-chi/chi/v5"
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
		log.Printf("failed to get products: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	SuccessResponse(w, http.StatusOK, ToProductResponses(products))
}

func (ih *InventoryHandler) InsertProduct(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to add products")
		return
	}
	var req InsertProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
		if errors.Is(err, service.ErrInvalidInput) {
			ErrorResponse(w, http.StatusBadRequest, "invalid input")
			return
		}
		log.Printf("failed to insert product: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	SuccessResponse(w, http.StatusCreated, ToProductResponse(productCreated))
}

func (ih *InventoryHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id param")
		return
	}
	if !isAdmin(r) {
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
	input := service.UpdateQuantityInput{
		ID:       id,
		Quantity: req.Quantity,
	}
	productModel, err := ih.serv.UpdateQuantity(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			ErrorResponse(w, http.StatusBadRequest, "invalid input")
			return
		}
		log.Printf("failed to update quantity: %s", err.Error())
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	SuccessResponse(w, http.StatusOK, ToProductResponse(productModel))
}

func (ih *InventoryHandler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id param")
		return
	}
	if !isAdmin(r) {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to change product price")
		return
	}
	var req struct {
		PriceCents int64 `json:"price_cents"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid req body")
		return
	}
	input := service.UpdateProductInput{
		ID:         id,
		PriceCents: req.PriceCents,
	}
	priceModel, err := ih.serv.UpdatePrice(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			ErrorResponse(w, http.StatusBadRequest, "invalid input")
			return
		}
		log.Printf("failed to update price: %s", err.Error())
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	SuccessResponse(w, http.StatusOK, ToProductResponse(priceModel))
}

func (h *InventoryHandler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	var req QuoteProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid input")
		return
	}
	input := service.GetQuotesInput{
		IDs: req.IDs,
	}
	productQuotes, err := h.serv.GetQuotes(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			ErrorResponse(w, http.StatusBadRequest, "invalid input")
			return
		}
		log.Printf("failed to update price: %s", err.Error())
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	SuccessResponse(w, http.StatusOK, ToQuotesProductReponses(productQuotes))
}

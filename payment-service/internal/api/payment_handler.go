package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	_ "github.com/Daty26/order-system/payment-service/internal/model"
	"github.com/Daty26/order-system/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}
type PaymentRequest struct {
	OrderID int     `json:"orderId"`
	Amount  float64 `json:"amount"`
}

func NewRepoPyament(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Process a new payment for a given order
// @Accept  json
// @Produce  json
// @Param payment body PaymentRequest true "Payment request"
// @Success 201 {object} model.Payment
// @Failure 400 {string} string "Invalid request"
// @Router /payments [post]
func (s *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID int     `json:"orderId"`
		Amount  float64 `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	payment, err := s.paymentService.ProcessPayment(req.OrderID, req.Amount)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	SuccessPayment(w, http.StatusCreated, payment)
}

func (s *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid id type")
		return
	}
	payment, err := s.paymentService.GetPaymentByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(w, http.StatusNotFound, "Can't find payment with specified id ")
		return
	}
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Can't fetch payment")
		return
	}
	SuccessPayment(w, http.StatusOK, payment)
}

// GetPayments godoc
// @Summary Get all payments
// @Description Retrieve a list of all processed payments
// @Produce  json
// @Success 200 {array} model.Payment
// @Failure 400 {string} string "Couldn't fetch payments"
// @Router /payments [get]
func (s *PaymentHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := s.paymentService.GetAllPayments()
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Couldn't fetch orders")
		return
	}
	SuccessPayment(w, http.StatusOK, payments)
}

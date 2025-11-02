package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Daty26/order-system/payment-service/internal/model"
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
		OrderID int `json:"orderId"`
		Amount  int `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid Request: "+err.Error())
		return
	}
	payment, err := s.paymentService.ProcessPayment(req.OrderID, req.Amount)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	SuccessPayment(w, http.StatusCreated, payment)
}

// UpdatePayment godoc
// @Summary Update an existing payment
// @Description Update status and amount of a payment by ID
// @Param id path int true "Payment ID"
// @Accept json
// @Produce json
// @Success 200 {object} model.Payment
// @Failure 400 {string} string "Invalid id or request body"
// @Failure 404 {string} string "Payment not found"
// @Failure 500 {string} string "Couldn't update payment"
// @Router /payments/{id} [put]
func (s *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid id type")
		return
	}
	var req struct {
		Status model.PaymentStatus `json:"status"`
		Amount float64             `json:"amount"`
	}
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	payment, err := s.paymentService.UpdatePayment(id, req.Status, req.Amount)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't update the payment: "+err.Error())
		return
	}
	SuccessPayment(w, http.StatusOK, payment)
}

// GetPaymentByID godoc
// @Summary Get payment by ID
// @Description Retrieve a single payment by its ID
// @Param id path int true "Payment ID"
// @Produce json
// @Success 200 {object} model.Payment
// @Failure 400 {string} string "Invalid id type"
// @Failure 404 {string} string "Payment not found"
// @Failure 500 {string} string "Could not fetch payment"
// @Router /payments/{id} [get]
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
		ErrorResponse(w, http.StatusInternalServerError, "Can't fetch payment: "+err.Error())
		return
	}
	SuccessPayment(w, http.StatusOK, payment)
}

// DeletePayment godoc
// @Summary Delete a payment
// @Description Delete payment with specified ID
// @Param id path int true "Payment ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid id type"
// @Failure 404 {string} string "Payment not found"
// @Failure 500 {string} string "Couldn't delete payment"
// @Router /payments/{id} [delete]
func (s *PaymentHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid id type")
		return
	}
	err = s.paymentService.DeletePayment(id)
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(w, http.StatusNotFound, "no payment with such id")
		return
	}
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't delete the payment")
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
		ErrorResponse(w, http.StatusBadRequest, "Couldn't fetch orders: "+err.Error())
		return
	}
	SuccessPayment(w, http.StatusOK, payments)
}

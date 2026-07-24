package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Daty26/order-system/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	logger         *slog.Logger
}

func NewPaymentHandler(paymentService *service.PaymentService, logger *slog.Logger) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, logger: logger}
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
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	userID := int(r.Context().Value("user_id").(float64))
	var req ProcessPaymentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}
	input := service.ProcessPaymentInput{
		OrderID:    req.OrderID,
		UserID:     userID,
		AuthHeader: r.Header.Get("Authorization"),
	}
	payment, err := h.paymentService.ProcessPayment(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			ErrorResponse(w, http.StatusBadRequest, "invalid input")
		default:
			h.logger.ErrorContext(r.Context(),
				"failed to process payment",
				"error", err,
				"order_id", input.OrderID,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
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
	var req UpdatePaymentRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid payment request")
		return
	}
	input := service.UpdatePaymentInput{
		ID:     id,
		Status: req.Status,
	}
	payment, err := s.paymentService.UpdatePayment(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			ErrorResponse(w, http.StatusBadRequest, "invalid payment request")
			s.logger.WarnContext(r.Context(), "incorrect payment request", "error", err)
		case errors.Is(err, sql.ErrNoRows):
			ErrorResponse(w, http.StatusNotFound, "payment not found")
		default:
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
			s.logger.ErrorContext(r.Context(), "failed to update the payment", "error", err, "payment_id", id)
		}
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
		ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	payment, err := s.paymentService.GetPaymentByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ErrorResponse(w, http.StatusNotFound, "not found")
		case errors.Is(err, service.ErrInvalidInput):
			ErrorResponse(w, http.StatusBadRequest, "invalid payment request")
		default:
			s.logger.ErrorContext(r.Context(), "failed fetch payment", "error", err, "payment_id", id)
			ErrorResponse(w, http.StatusInternalServerError, "someting went wrong")
		}
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
	err = s.paymentService.DeletePayment(r.Context(), id)
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
	role, ok := r.Context().Value("role").(string)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
	userIDRaw, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
	limit, offset, ok := parsePagination(r)
	if !ok {
		ErrorResponse(w, http.StatusBadRequest, "invalid pagination query param")
		return
	}
	if role == "ADMIN" {
		payments, err := s.paymentService.GetAllPayments(r.Context(), limit, offset)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidInput):
				ErrorResponse(w, http.StatusBadRequest, "invalid payment request")
			default:
				s.logger.ErrorContext(r.Context(), "failed to get payments", "error", err, "user_id", int(userIDRaw))
				ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
			}
			return
		}
		SuccessPayment(w, http.StatusOK, payments)
		return
	}
	input := service.GetAllByUserIDInput{
		ID:     int(userIDRaw),
		Limit:  limit,
		Offset: offset,
	}
	payments, err := s.paymentService.GetAllByUserId(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			ErrorResponse(w, http.StatusBadRequest, "invalid request")
		default:
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	SuccessPayment(w, http.StatusOK, payments)
}

package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Daty26/order-system/payment-service/internal/service"
)

func HandleErrors(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	err error,
	msg string,
	args ...any,
) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		ErrorResponse(w, http.StatusBadRequest, "invalid payment request")
	case errors.Is(err, service.ErrPaymentAlreadyExists):
		ErrorResponse(w, http.StatusConflict, "payment already exists")
	default:
		logger.ErrorContext(
			r.Context(),
			msg,
			append([]any{"error", err}, args...)...,
		)
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
	}
}

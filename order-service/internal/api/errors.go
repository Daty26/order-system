package api

import (
	"database/sql"
	"errors"
	"github.com/Daty26/order-system/order-service/internal/service"
	"log/slog"
	"net/http"
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
	case errors.Is(err, service.ErrInvalidRequest):
		ErrorResponse(w, http.StatusBadRequest, "invalid order request")
	case errors.Is(err, service.ErrForbiddenOrder):
		ErrorResponse(w, http.StatusForbidden, "failed to modify order")
	case errors.Is(err, service.ErrOrderCannotBeCanceled):
		ErrorResponse(w, http.StatusBadRequest, "failed to cancel the order")
	case errors.Is(err, service.ErrProductNotFound):
		ErrorResponse(w, http.StatusNotFound, "products not found")
	case errors.Is(err, sql.ErrNoRows):
		ErrorResponse(w, http.StatusNotFound, "order not found")
	default:
		logger.ErrorContext(r.Context(), msg, append([]any{"error", err}, args...)...)
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
	}

}

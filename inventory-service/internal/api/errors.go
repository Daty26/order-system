package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Daty26/order-system/inventory-service/internal/service"
)

func HandleErrors(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	err error,
	msg string,
	attrs ...any,
) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		ErrorResponse(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, service.ErrInsufficientStock):
		ErrorResponse(w, http.StatusConflict, "stock amount is insufficient")
	case errors.Is(err, sql.ErrNoRows):
		ErrorResponse(w, http.StatusNotFound, "not found")
	default:
		logger.ErrorContext(r.Context(), msg, append([]any{"error", err}, attrs...)...)
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
	}
}

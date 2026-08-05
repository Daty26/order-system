package api

import (
	"database/sql"
	"errors"
	"github.com/Daty26/order-system/notification-service/internal/service"
	"log/slog"
	"net/http"
)

func HandleError(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	err error,
	msg string,
	attrs ...any,
) {
	switch {
	case errors.Is(err, service.ErrInvalidID):
		ErrorResponse(w, http.StatusBadRequest, "invalid id")
	case errors.Is(err, service.ErrInvalidStatus):
		ErrorResponse(w, http.StatusBadRequest, "invalid status")
	case errors.Is(err, service.ErrInvalidMessage):
		ErrorResponse(w, http.StatusBadRequest, "invalid message")
	case errors.Is(err, service.ErrNotificationAlreadyExists):
		ErrorResponse(w, http.StatusConflict, "notification already exists")
	case errors.Is(err, sql.ErrNoRows):
		ErrorResponse(w, http.StatusNotFound, "not found")
	default:
		logger.ErrorContext(r.Context(), msg, append([]any{"error", err}, attrs...)...)
		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
	}

}

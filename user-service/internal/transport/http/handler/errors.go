package transport_http_handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Daty26/order-system/user-service/internal/service"
	transport_http_response "github.com/Daty26/order-system/user-service/internal/transport/http/response"
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
	case errors.Is(err, service.ErrIncorrectID):
		transport_http_response.ErrorJSON(w, http.StatusBadRequest, "invalid request id")
	case errors.Is(err, service.ErrInvalidCredentials):
		transport_http_response.ErrorJSON(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, service.ErrInvalidUserInput):
		transport_http_response.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
	case errors.Is(err, service.ErrNotFound):
		transport_http_response.ErrorJSON(w, http.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrUserAlreadyExists):
		transport_http_response.ErrorJSON(w, http.StatusConflict, "user already exists")
	default:
		logger.ErrorContext(
			r.Context(),
			msg,
			append([]any{"error", err}, args...)...,
		)
		transport_http_response.ErrorJSON(w, http.StatusInternalServerError, "something went wrong")
	}

}

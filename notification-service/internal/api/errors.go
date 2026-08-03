package api

import (
	"log/slog"
	"net/http"
)

// TODO implement HandleError and use it in all controllers
func HandleError(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	err error,
	msg string,
	attrs ...any,
) {

}

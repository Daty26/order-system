package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Daty26/order-system/notification-service/internal/model"
	"github.com/Daty26/order-system/notification-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct {
	s      *service.NotificationService
	logger *slog.Logger
}

func NewNotificationHardler(service *service.NotificationService, logger *slog.Logger) *NotificationHandler {
	return &NotificationHandler{s: service, logger: logger}
}

func (h *NotificationHandler) InsertNotification(w http.ResponseWriter, r *http.Request) {
	userIDRaw, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID := int(userIDRaw)
	var req InsertNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	input := service.InsertInput{
		OrderID:   req.OrderID,
		PaymentID: req.PaymentID,
		Status:    model.NotificationStatus(req.Status),
		Message:   req.Message,
		UserID:    userID,
	}
	notification, err := h.s.Insert(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidStatus):
			ErrorResponse(w, http.StatusBadRequest, "invalid status")
		case errors.Is(err, service.ErrInvalidMessage):
			ErrorResponse(w, http.StatusBadRequest, "invalid message")
		case errors.Is(err, service.ErrInvalidID):
			ErrorResponse(w, http.StatusBadRequest, "invalid id")
		default:
			h.logger.ErrorContext(
				r.Context(),
				"failed to create notification",
				"error", err,
				"user_id", userID,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	SuccessResp(w, http.StatusCreated, notification)
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userIDRaw, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID := int(userIDRaw)
	limit, offset, ok := parsePagination(r)
	if !ok {
		ErrorResponse(w, http.StatusBadRequest, "invalid pagination")
		return
	}
	if role == "ADMIN" {
		notifications, err := h.s.GetAll(r.Context(), limit, offset)
		if err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"failed to retrieve notifications",
				"err", err,
				"user_id", userID,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		SuccessResp(w, http.StatusOK, notifications)
		return
	}
	input := service.GetAllByUserIDInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
	notifications, err := h.s.GetAllByUserID(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			ErrorResponse(w, http.StatusBadRequest, "invalid id")
		default:
			h.logger.ErrorContext(
				r.Context(),
				"failed to retrieve notifications",
				"error", err,
				"user_id", userID,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	SuccessResp(w, http.StatusOK, notifications)
	return
}

func (h *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	notificationModel, err := h.s.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ErrorResponse(w, http.StatusNotFound, "notification not found")
		default:
			h.logger.ErrorContext(
				r.Context(),
				"failed to retrieve notification",
				"error", err,
				"notification_id", id,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	SuccessResp(w, http.StatusOK, notificationModel)
}

func (h *NotificationHandler) GetNotificationsByStatus(w http.ResponseWriter, r *http.Request) {
	userIDRaw, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusBadRequest, "unauthorized")
		return
	}
	userID := int(userIDRaw)
	limit, offset, ok := parsePagination(r)
	if !ok {
		ErrorResponse(w, http.StatusBadRequest, "invalid pagination params")
		return
	}

	input := service.GetByStatusInput{
		Status: model.NotificationStatus(chi.URLParam(r, "status")),
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
	notifications, err := h.s.GetByStatus(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidStatus):
			ErrorResponse(w, http.StatusBadRequest, "invalid status")
		case errors.Is(err, service.ErrInvalidID):
			ErrorResponse(w, http.StatusBadRequest, "invalid id")
		default:
			h.logger.ErrorContext(
				r.Context(),
				"failed to retrieve notification by status",
				"error", err,
				"user_id", userID,
				"status", input.Status,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	SuccessResp(w, http.StatusOK, notifications)
}

func (h *NotificationHandler) UpdateNotificationStatus(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	if role != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "You don't have permission to update notification")
		return
	}
	var req UpdateNotificationByStatus
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	input := service.UpdateStatusInput{
		Status: model.NotificationStatus(req.Status),
		ID:     id,
	}
	notification, err := h.s.UpdateStatus(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			ErrorResponse(w, http.StatusBadRequest, "invalid id")
		case errors.Is(err, service.ErrInvalidStatus):
			ErrorResponse(w, http.StatusBadRequest, "invalid status")
		default:
			h.logger.ErrorContext(
				r.Context(),
				"failed to update notification",
				"error", err,
				"status", input.Status,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	SuccessResp(w, http.StatusOK, notification)
}

func (h *NotificationHandler) DeleteNotificationByID(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if role != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "You don't have permission to update notification")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err = h.s.DeleteByID(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			ErrorResponse(w, http.StatusBadRequest, "invalid id")
		case errors.Is(err, sql.ErrNoRows):
			ErrorResponse(w, http.StatusNotFound, "notification not found")
		default:
			h.logger.ErrorContext(
				r.Context(),
				"failed to delete notification",
				"error", err,
				"notification_id", id,
			)
			ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}

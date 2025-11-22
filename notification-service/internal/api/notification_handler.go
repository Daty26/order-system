package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Daty26/order-system/notification-service/internal/model"
	"github.com/Daty26/order-system/notification-service/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type NotificationHandler struct {
	sv *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{sv: service}
}
func (nh *NotificationHandler) InsertNotification(w http.ResponseWriter, r *http.Request) {
	userIdFloat := r.Context().Value("user_id").(float64)
	userId := int(userIdFloat)

	var req struct {
		OrderID   int                      `json:"order_id"`
		PaymentID int                      `json:"payment_id"`
		Status    model.NotificationStatus `json:"status"`
		Message   string                   `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	notification, err := nh.sv.Insert(req.OrderID, req.PaymentID, req.Status, req.Message, userId)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't insert new payment: "+err.Error())
		return
	}
	SuccessResp(w, http.StatusCreated, notification)
}

func (nh *NotificationHandler) GetAllNotificationsByUserId(w http.ResponseWriter, r *http.Request) {
	userIdFloat := r.Context().Value("user_id").(float64)
	userId := int(userIdFloat)
	notifications, err := nh.sv.GetAllByUserID(userId)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't retrieve notifications: "+err.Error())
		return
	}
	SuccessResp(w, http.StatusOK, notifications)
}

func (nh *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "incorrect id format"+err.Error())
		return
	}
	notificationByID, err := nh.sv.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(w, http.StatusNotFound, "Couldn't find notification with such id")
		return
	}
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResp(w, http.StatusOK, notificationByID)
}
func (nh *NotificationHandler) GetNotificationsByStatus(w http.ResponseWriter, r *http.Request) {
	userIdFloat := r.Context().Value("user_id").(float64)
	userId := int(userIdFloat)
	status := model.NotificationStatus(chi.URLParam(r, "status"))
	notifications, err := nh.sv.GetByStatus(status, userId)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResp(w, http.StatusOK, notifications)
}
func (nh *NotificationHandler) UpdateNotificationStatusByID(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value("role")
	if userRole != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "You don't have permission to update notification")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "incorrect id: "+err.Error())
		return
	}
	var req struct {
		Status model.NotificationStatus `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	notification, err := nh.sv.UpdateStatusByID(id, req.Status)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't update the status: "+err.Error())
		return
	}
	SuccessResp(w, http.StatusOK, notification)
}
func (nh *NotificationHandler) DeleteNotificationByID(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value("role")
	if userRole != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "You don't have permission to update notification")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "incorrect id: "+err.Error())
		return
	}
	if err = nh.sv.DeleteByID(id); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal error: "+err.Error())
	}
	SuccessResp(w, http.StatusOK, nil)
}

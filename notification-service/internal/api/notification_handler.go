package api

import (
	"encoding/json"
	"github.com/Daty26/order-system/notification-service/internal/model"
	"github.com/Daty26/order-system/notification-service/internal/service"
	"net/http"
)

type NotificationHandler struct {
	sv *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{sv: service}
}
func (nh *NotificationHandler) InsertNotification(w http.ResponseWriter, r *http.Request) {

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
	notification, err := nh.sv.Insert(req.OrderID, req.PaymentID, req.Status, req.Message)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't insert new payment: "+err.Error())
		return
	}
	SuccessPayment(w, http.StatusCreated, notification)
	//finish later!
}

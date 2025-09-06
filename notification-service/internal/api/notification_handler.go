package api

import (
	"encoding/json"
	"github.com/Daty26/order-system/notification-service/internal/service"
	"net/http"
)

type NotificationHandler struct {
	sv *service.NotificationService
}

func (nh *NotificationHandler) InsertNotification(w http.ResponseWriter, r *http.Request) {
	var req struct {
	}
	json.NewDecoder(r.Body).Decode(&req)
	nh.sv.Insert()
	//finish later!
}

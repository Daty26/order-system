package api

import (
	"net/http"

	"github.com/Daty26/order-system/inventory-service/internal/service"
)

func actorFromContext(r *http.Request) (service.Actor, bool) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		return service.Actor{}, false
	}
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		return service.Actor{}, false
	}
	return service.NewActor(int(userID), role), true
}

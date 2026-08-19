package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Daty26/order-system/payment-service/internal/service"
)

type OrderCreatedHandler struct {
	service *service.PaymentService
	logger  *slog.Logger
}

func NewOrderCreatedHandler(service *service.PaymentService, logger *slog.Logger) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		service: service,
		logger:  logger,
	}
}

func (h *OrderCreatedHandler) Handle(value []byte) {
	var event service.OrderCreatedEvent
	if err := json.Unmarshal(value, &event); err != nil {
		h.logger.Error("invalid order.created event", "error", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := h.service.ProcessOrderCreated(ctx, event); err != nil {
		h.logger.Error("failed to process order.created", "error", err, "order_id", event.OrderID)
		return
	}
	h.logger.Info("processed order.created", "order_id", event.OrderID)
}

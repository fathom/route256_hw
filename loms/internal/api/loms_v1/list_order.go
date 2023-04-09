package loms_v1

import (
	"context"
	"fmt"
	"route256/loms/internal/logger"
	desc "route256/loms/pkg/loms_v1"
)

// listOrder
// Показывает информацию по заказу

func (h *Handlers) ListOrder(ctx context.Context, request *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	logger.Debug(fmt.Sprintf("listOrder: %+v", request))

	order, orderItems, err := h.businessLogic.ListOrder(ctx, request.GetOrderId())
	if err != nil {
		return nil, err
	}

	var items []*desc.OrderItem

	for _, item := range orderItems {
		items = append(items, &desc.OrderItem{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	response := &desc.ListOrderResponse{
		Status: string(order.Status),
		User:   order.UserID,
		Items:  items,
	}

	return response, nil
}

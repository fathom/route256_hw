package loms_v1

import (
	"context"
	"log"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

//listOrder
//Показывает информацию по заказу

func (h *Handlers) ListOrder(_ context.Context, request *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	log.Printf("listOrder: %+v", request)

	var items []*desc.OrderItem

	items = append(items, &desc.OrderItem{
		Sku:   1148162,
		Count: 2,
	}, &desc.OrderItem{
		Sku:   1076963,
		Count: 1,
	})

	response := &desc.ListOrderResponse{
		Status: string(domain.AwaitingPayment),
		User:   100,
		Items:  items,
	}

	_ = domain.NewStatus
	_ = domain.Cancelled
	_ = domain.Failed
	_ = domain.Payed

	return response, nil
}

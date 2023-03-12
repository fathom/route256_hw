package loms_v1

import (
	"context"
	"log"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"
)

// createOrder
// Создает новый заказ для пользователя из списка переданных товаров.
// Товары при этом нужно зарезервировать на складе.

func (h *Handlers) CreateOrder(ctx context.Context, request *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	log.Printf("createOrder: %+v", request)

	orderID, err := h.businessLogic.CreateOrder(ctx, request.GetUser(), converter.ToOrderItemListDomain(request.GetItems()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateOrderResponse{
		OrderId: orderID,
	}, nil
}

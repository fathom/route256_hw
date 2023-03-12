package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// cancelOrder
// Отменяет заказ, снимает резерв со всех товаров в заказе.

func (h *Handlers) CancelOrder(ctx context.Context, request *desc.CancelOrderRequest) (*emptypb.Empty, error) {
	log.Printf("cancelOrder: %+v", request)

	err := h.businessLogic.CancelOrder(ctx, request.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

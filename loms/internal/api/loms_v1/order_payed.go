package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// orderPayed
// Помечает заказ оплаченным. Зарезервированные товары должны перейти в статус купленных.

func (h *Handlers) OrderPayed(ctx context.Context, request *desc.OrderPayedRequest) (*emptypb.Empty, error) {
	log.Printf("orderPayed: %+v", request)

	err := h.businessLogic.OrderPayed(ctx, request.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

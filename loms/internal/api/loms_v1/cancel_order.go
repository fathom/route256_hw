package loms_v1

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	desc "route256/loms/pkg/loms_v1"
)

//cancelOrder
//Отменяет заказ, снимает резерв со всех товаров в заказе.

func (h *Handlers) CancelOrder(_ context.Context, request *desc.CancelOrderRequest) (*emptypb.Empty, error) {
	log.Printf("cancelOrder: %+v", request)

	return &emptypb.Empty{}, nil
}

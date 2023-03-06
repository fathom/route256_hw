package checkout_v1

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

//purchase
//Оформить заказ по всем товарам корзины. Вызывает createOrder у LOMS.

func (h *Handlers) Purchase(ctx context.Context, request *desc.PurchaseRequest) (*emptypb.Empty, error) {
	log.Printf("purchase: %+v", request)

	err := h.businessLogic.Purchase(ctx, request.GetUser())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

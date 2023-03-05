package checkout_v1

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

//deleteFromCart
//Удалить товар из корзины определенного пользователя.

func (h *Handlers) DeleteFromCart(_ context.Context, request *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	log.Printf("deleteFromCart: %+v", request)

	return &emptypb.Empty{}, nil
}
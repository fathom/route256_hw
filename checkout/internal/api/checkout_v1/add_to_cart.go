package checkout_v1

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	desc "route256/checkout/pkg/checkout_v1"
)

//addToCart
//Добавить товар в корзину определенного пользователя. При этом надо
//проверить наличие товара через LOMS.stocks

func (h *Handlers) AddToCart(ctx context.Context, request *desc.AddToCartRequest) (*emptypb.Empty, error) {
	log.Printf("addToCart: %+v", request)

	err := h.businessLogic.AddToCart(ctx, request.GetUser(), request.GetSku(), request.GetCount())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

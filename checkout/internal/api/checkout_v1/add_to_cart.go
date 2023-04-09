package checkout_v1

import (
	"context"
	"fmt"
	"route256/checkout/internal/logger"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// addToCart
// Добавить товар в корзину определенного пользователя. При этом надо
// проверить наличие товара через LOMS.stocks

func (h *Handlers) AddToCart(ctx context.Context, request *desc.AddToCartRequest) (*emptypb.Empty, error) {
	logger.Debug(fmt.Sprintf("addToCart: %+v", request))

	err := h.businessLogic.AddToCart(ctx, request.GetUser(), request.GetSku(), request.GetCount())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

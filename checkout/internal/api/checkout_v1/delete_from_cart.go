package checkout_v1

import (
	"context"
	"fmt"
	"route256/checkout/internal/logger"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// deleteFromCart
// Удалить товар из корзины определенного пользователя.

func (h *Handlers) DeleteFromCart(ctx context.Context, request *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	logger.Debug(fmt.Sprintf("deleteFromCart: %+v", request))

	err := h.businessLogic.DeleteFromCart(ctx, request.User, request.Sku, request.Count)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

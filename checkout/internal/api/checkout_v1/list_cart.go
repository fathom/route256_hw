package checkout_v1

import (
	"context"
	"fmt"
	"route256/checkout/internal/logger"
	desc "route256/checkout/pkg/checkout_v1"
)

// listCart
// Показать список товаров в корзине с именами и ценами (их надо в реальном
// времени получать из ProductService)

func (h *Handlers) ListCart(ctx context.Context, request *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	logger.Debug(fmt.Sprintf("listCart: %+v", request))

	cartItems, err := h.businessLogic.ListCart(ctx, request.GetUser())
	if err != nil {
		return nil, err
	}

	response := &desc.ListCartResponse{}

	for _, item := range cartItems {
		response.Items = append(response.Items, &desc.Product{
			Sku:   item.Sku,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
		response.TotalPrice += item.Price
	}

	return response, nil
}

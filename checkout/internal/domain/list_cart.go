package domain

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

type CartItem struct {
	Scu   uint32
	Count uint16
	Name  string
	Price uint32
}

func (m *Domain) ListCart(ctx context.Context, user int64) ([]CartItem, error) {
	log.Printf("listCart for user: %+v", user)

	userCart := []CartItem{
		{1076963, 1, "", 0},
		{1148162, 1, "", 0},
		{1625903, 1, "", 0},
		{2618151, 1, "", 0},
	}

	for i := range userCart {
		name, price, err := m.productService.GetProduct(ctx, userCart[i].Scu)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong sku")
		}

		userCart[i].Name = name
		userCart[i].Price = price
	}

	return userCart, nil
}

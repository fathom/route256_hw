package domain

import (
	"context"
	"log"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

func (d *domain) ListCart(ctx context.Context, userID int64) ([]model.CartItem, error) {
	log.Printf("listCart for user: %+v", userID)

	userCart, err := d.cartRepository.ListCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range userCart {
		name, price, err := d.productService.GetProduct(ctx, userCart[i].Sku)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong sku")
		}

		userCart[i].Name = name
		userCart[i].Price = price
	}

	return userCart, nil
}

package domain

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrNotFound           = errors.New("not found")
)

func (d *domain) AddToCart(ctx context.Context, userID int64, sku uint32, count uint32) error {

	stocks, err := d.lomsService.Stocks(ctx, sku)
	if err != nil {
		return errors.WithMessage(err, "checking stocks")
	}

	log.Printf("stocks reseived: %+v", stocks)

	cartItem, err := d.cartRepository.GetCartItem(ctx, userID, sku)
	if err != nil && !errors.As(err, &ErrNotFound) {
		log.Printf("GetCartItem: %+v", err)
		return err
	}

	counter := int64(count + cartItem.Count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			err = d.cartRepository.AddToCart(ctx, userID, sku, count+cartItem.Count)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return ErrInsufficientStocks
}

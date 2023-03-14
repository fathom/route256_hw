package domain

import (
	"context"
)

func (d *domain) DeleteFromCart(ctx context.Context, userID int64, sku uint32, count uint32) error {

	cartItem, err := d.cartRepository.GetCartItem(ctx, userID, sku)
	if err != nil {
		return err
	}

	if cartItem.Count <= count {
		err = d.cartRepository.DeleteCart(ctx, userID, sku)
		if err != nil {
			return err
		}
	} else {
		newCount := cartItem.Count - count
		err = d.cartRepository.UpdateCountCart(ctx, userID, sku, newCount)
		if err != nil {
			return err
		}
	}

	return nil
}

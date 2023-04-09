package domain

import (
	"context"
	"fmt"
	"route256/checkout/internal/logger"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

func (d *domain) Purchase(ctx context.Context, userID int64) error {
	logger.Debug(fmt.Sprintf("CreateOrder for user: %+v", userID))

	userCart, err := d.cartRepository.ListCart(ctx, userID)
	if err != nil {
		return err
	}

	var items []*model.OrderItem

	for _, item := range userCart {
		items = append(items, &model.OrderItem{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	orderID, err := d.lomsService.CreateOrder(ctx, userID, items)
	if err != nil {
		return errors.WithMessage(err, "failed create order")
	}
	logger.Debug(fmt.Sprintf("Created Order: %+v", orderID))

	err = d.cartRepository.DeleteUserCart(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

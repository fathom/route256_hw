package domain

import (
	"context"
	"route256/loms/internal/model"
)

func (d *Domain) ListOrder(ctx context.Context, orderId int64) (model.Order, []model.OrderItem, error) {
	order, err := d.OrdersRepository.GetOrder(ctx, orderId)
	if err != nil {
		return model.Order{}, nil, err
	}

	items, err := d.OrderItemsRepository.GetItemsByOrderId(ctx, orderId)
	if err != nil {
		return model.Order{}, nil, err
	}

	return order, items, nil
}

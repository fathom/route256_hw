package domain

import (
	"context"
	"errors"
	"log"
	"route256/loms/internal/model"
)

var (
	ErrCreateOrderFailed = errors.New("create order failed")
)

func (d *Domain) CreateOrder(ctx context.Context, user int64, items []*model.OrderItem) (int64, error) {

	newOrder := model.Order{
		Status: model.NewStatus,
		UserID: user,
	}

	var orderId int64

	err := d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {

		var err error

		orderId, err = d.OrdersRepository.CreateOrder(ctxTX, newOrder)
		if err != nil {
			return err
		}

		for _, item := range items {
			item.OrderId = orderId
			err = d.OrderItemsRepository.CreateItem(ctxTX, *item)
			if err != nil {
				return err
			}
		}

		// todo резервация товаров на складе

		return nil
	})

	if err != nil {
		log.Println("create order failed", err)
		return 0, ErrCreateOrderFailed
	}

	return orderId, nil
}

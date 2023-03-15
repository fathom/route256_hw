package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"route256/loms/internal/model"
	"time"
)

var (
	ErrCreateOrderFailed = errors.New("create order failed")
	ErrReservationFailed = errors.New("reservation failed")
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
			item.OrderID = orderId

			err = d.OrderItemsRepository.CreateItem(ctxTX, *item)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Println("create order failed", err)
		return 0, ErrCreateOrderFailed
	}

	err = d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {

		expiredAt := time.Now().Add(time.Minute * 10)

		for _, item := range items {
			stocks, err := d.WarehouseRepository.GetStocksBySku(ctxTX, item.Sku)
			if err != nil {
				return err
			}
			for _, w := range stocks {
				count := w.Count - w.Reservation
				if count > 0 {
					reserved := item.Count
					if count < reserved {
						reserved = count
					}

					item.Count -= reserved

					err := d.WarehouseRepository.AddReservation(ctxTX, item.Sku, w.WarehouseID, orderId, reserved, expiredAt)
					if err != nil {
						return err
					}
				}

				if item.Count <= 0 {
					break
				}
			}

			if item.Count > 0 {
				return fmt.Errorf("reservation failed item %v", item.Sku)
			}
		}

		return nil
	})

	if err != nil {
		log.Println("reservation failed", err)
		err := d.OrdersRepository.UpdateStatusOrder(ctx, orderId, model.Failed)
		if err != nil {
			return 0, err
		}
		return 0, ErrReservationFailed
	}

	return orderId, nil
}

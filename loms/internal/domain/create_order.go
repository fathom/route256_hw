package domain

import (
	"context"
	"errors"
	"fmt"
	"route256/loms/internal/logger"
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
		logger.Debug(fmt.Sprintf("create order failed: %+v", err))
		return 0, ErrCreateOrderFailed
	}

	err = d.OrderStatusSender.SendOrderStatus(orderId, newOrder.Status)
	if err != nil {
		logger.Debug(fmt.Sprintf("OrderStatusSender: %+v", err))
		return 0, err
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
		logger.Debug(fmt.Sprintf("reservation failed: %+v", err))
		err := d.OrdersRepository.UpdateStatusOrder(ctx, orderId, model.Failed)
		if err != nil {
			return 0, err
		}

		err = d.OrderStatusSender.SendOrderStatus(orderId, model.Failed)
		if err != nil {
			logger.Debug(fmt.Sprintf("OrderStatusSender: %+v", err))
			return 0, err
		}

		return 0, ErrReservationFailed
	}

	err = d.OrdersRepository.UpdateStatusOrder(ctx, orderId, model.AwaitingPayment)
	if err != nil {
		return 0, err
	}

	// Отправляем отложенную задачу на удаление заказа
	d.DeleteReservationWorker.AddDelayJob(model.JobDeleteReservation{
		OrderId: orderId,
	})

	err = d.OrderStatusSender.SendOrderStatus(orderId, model.AwaitingPayment)
	if err != nil {
		logger.Debug(fmt.Sprintf("OrderStatusSender: %+v", err))
		return 0, err
	}

	return orderId, nil
}

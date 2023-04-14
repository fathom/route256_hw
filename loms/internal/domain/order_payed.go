package domain

import (
	"context"
	"errors"
	"fmt"
	"route256/loms/internal/logger"
	"route256/loms/internal/model"
)

var (
	ErrPayOrderFailed = errors.New("pay order failed")
)

func (d *Domain) OrderPayed(ctx context.Context, orderId int64) error {

	err := d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		err := d.OrdersRepository.UpdateStatusOrder(ctxTX, orderId, model.Payed)
		if err != nil {
			return err
		}
		logger.Debug(fmt.Sprintf("order %v mark as %v", orderId, model.Payed))

		reservations, err := d.WarehouseRepository.GetReservationByOrderId(ctxTX, orderId)
		if err != nil {
			return err
		}

		for _, r := range reservations {
			stocks, err := d.WarehouseRepository.GetStocks(ctxTX, r.Sku, r.WarehouseID)
			if err != nil {
				return err
			}
			err = d.WarehouseRepository.UpdateWarehouse(ctxTX, r.Sku, r.WarehouseID, stocks.Count-r.Count)
			if err != nil {
				return err
			}
		}

		err = d.WarehouseRepository.DeleteReservation(ctxTX, orderId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Debug(fmt.Sprintf("pay order failed: %+v", err))
		return ErrPayOrderFailed
	}

	err = d.OrderStatusSender.SendOrderStatus(orderId, model.Payed)
	if err != nil {
		logger.Debug(fmt.Sprintf("OrderStatusSender: %+v", err))
		return err
	}

	return nil
}

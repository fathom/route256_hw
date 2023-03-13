package domain

import (
	"context"
	"errors"
	"log"
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
		log.Printf("order %v mark as %v", orderId, model.Payed)

		reservations, err := d.WarehouseRepository.GetReservationByOrderId(ctxTX, orderId)
		if err != nil {
			return err
		}

		for _, r := range reservations {
			err := d.WarehouseRepository.UpdateWarehouse(ctxTX, r.Sku, r.WarehouseID, r.Count)
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
		log.Println("pay order failed", err)
		return ErrPayOrderFailed
	}

	return nil
}

package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/logger"
	"route256/loms/internal/model"
)

func (d *Domain) CancelOrder(ctx context.Context, orderId int64) error {
	err := d.OrdersRepository.UpdateStatusOrder(ctx, orderId, model.Cancelled)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("order %v mark as %v", orderId, model.Cancelled))

	err = d.WarehouseRepository.DeleteReservation(ctx, orderId)
	if err != nil {
		return err
	}

	err = d.OrderStatusSender.SendOrderStatus(orderId, model.Cancelled)
	if err != nil {
		logger.Debug(fmt.Sprintf("OrderStatusSender: %+v", err))
		return err
	}

	return nil
}

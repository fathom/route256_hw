package domain

import (
	"context"
	"log"
	"route256/loms/internal/model"
)

func (d *Domain) CancelOrder(ctx context.Context, orderId int64) error {
	err := d.OrdersRepository.UpdateStatusOrder(ctx, orderId, model.Cancelled)
	if err != nil {
		return err
	}
	log.Printf("order %v mark as %v", orderId, model.Cancelled)
	// todo снимает резерв со всех товаров в заказе

	return nil
}

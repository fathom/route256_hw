package domain

import (
	"context"
	"log"
	"route256/loms/internal/model"
)

func (d *Domain) OrderPayed(ctx context.Context, orderId int64) error {
	err := d.OrdersRepository.UpdateStatusOrder(ctx, orderId, model.Payed)
	if err != nil {
		return err
	}
	log.Printf("order %v mark as %v", orderId, model.Payed)

	// todo Зарезервированные товары должны перейти в статус купленных.

	return nil
}

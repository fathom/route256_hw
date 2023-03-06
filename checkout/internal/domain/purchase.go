package domain

import (
	"context"
	"log"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

func (s *domain) Purchase(ctx context.Context, user int64) error {
	log.Printf("CreateOrder for user: %+v", user)

	items := []*model.OrderItem{
		{Sku: 1076963, Count: 1},
		{Sku: 1148162, Count: 1},
		{Sku: 1625903, Count: 1},
		{Sku: 2618151, Count: 1},
	}

	orderID, err := s.lomsService.CreateOrder(ctx, user, items)
	if err != nil {
		return errors.WithMessage(err, "failed create order")
	}
	log.Printf("Created Order: %+v", orderID)

	return nil
}

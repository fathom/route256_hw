package domain

import (
	"context"
	"github.com/pkg/errors"
	"log"
)

func (m *Domain) Purchase(ctx context.Context, user int64) error {
	log.Printf("CreateOrder for user: %+v", user)

	items := []OrderItem{
		{1076963, 1},
		{1148162, 1},
		{1625903, 1},
		{2618151, 1},
	}

	orderID, err := m.lomsService.CreateOrder(ctx, user, items)
	if err != nil {
		return errors.WithMessage(err, "failed create order")
	}
	log.Printf("Created Order: %+v", orderID)

	return nil
}

package domain

import (
	"context"
)

func (d *Domain) CreateOrder(ctx context.Context, user int64, items []*OrderItem) (int64, error) {

	_ = ctx
	_ = user
	_ = items

	return 123456, nil
}

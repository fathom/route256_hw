package domain

import (
	"context"
	"route256/loms/internal/model"
)

type OrdersRepository interface {
	GetOrder(context.Context, int64) (model.Order, error)
	CreateOrder(context.Context, model.Order) (int64, error)
	UpdateStatusOrder(context.Context, int64, model.OrderStatus) error
}

type OrderItemsRepository interface {
	CreateItem(context.Context, model.OrderItem) error
	GetItemsByOrderId(context.Context, int64) ([]model.OrderItem, error)
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
	// todo RunSerializable()
	// todo RunCommitted()
	// todo RunUncommitted()
}

type Domain struct {
	TransactionManager
	OrdersRepository
	OrderItemsRepository
}

func New(
	TransactionManager TransactionManager,
	OrdersRepository OrdersRepository,
	OrderItemsRepository OrderItemsRepository,
) *Domain {
	return &Domain{
		TransactionManager,
		OrdersRepository,
		OrderItemsRepository,
	}
}

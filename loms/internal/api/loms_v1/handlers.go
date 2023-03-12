package loms_v1

import (
	"context"
	"route256/loms/internal/model"
	desc "route256/loms/pkg/loms_v1"
)

type BusinessLogic interface {
	CreateOrder(context.Context, int64, []*model.OrderItem) (int64, error)
	ListOrder(context.Context, int64) (model.Order, []model.OrderItem, error)
	CancelOrder(context.Context, int64) error
	OrderPayed(context.Context, int64) error
	Stocks(context.Context, uint32) ([]model.StockItem, error)
}
type Handlers struct {
	desc.UnimplementedLomsV1Server
	businessLogic BusinessLogic
}

func NewLomsV1(businessLogic BusinessLogic) *Handlers {
	return &Handlers{
		businessLogic: businessLogic,
	}
}

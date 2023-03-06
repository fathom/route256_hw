package loms_v1

import (
	"context"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

type BusinessLogic interface {
	CreateOrder(context.Context, int64, []*domain.OrderItem) (int64, error)
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

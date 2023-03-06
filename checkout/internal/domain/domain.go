package domain

import (
	"context"
	"route256/checkout/internal/model"
)

type LomsService interface {
	Stocks(ctx context.Context, sku uint32) ([]*model.Stock, error)
	CreateOrder(ctx context.Context, user int64, items []*model.OrderItem) (int64, error)
}

type ProductService interface {
	GetProduct(ctx context.Context, sku uint32) (string, uint32, error)
	ListSkus(ctx context.Context, startAfterSku, count uint32) ([]uint32, error)
}

type BusinessLogic interface {
	AddToCart(context.Context, int64, uint32, uint32) error
	ListCart(context.Context, int64) ([]CartItem, error)
	Purchase(context.Context, int64) error
}

var _ BusinessLogic = (*domain)(nil)

type domain struct {
	lomsService    LomsService
	productService ProductService
}

func New(lomsService LomsService, productService ProductService) *domain {
	return &domain{
		lomsService:    lomsService,
		productService: productService,
	}
}

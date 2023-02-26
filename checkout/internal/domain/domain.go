package domain

import "context"

type LomsService interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (int64, error)
}

type ProductService interface {
	GetProduct(ctx context.Context, sku uint32) (string, uint32, error)
	ListSkus(ctx context.Context, startAfterSku, count uint32) ([]uint32, error)
}

type OrderItem struct {
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Domain struct {
	lomsService    LomsService
	productService ProductService
}

func New(lomsService LomsService, productService ProductService) *Domain {
	return &Domain{
		lomsService:    lomsService,
		productService: productService,
	}
}

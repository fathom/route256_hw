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
	ListCart(context.Context, int64) ([]model.CartItem, error)
	Purchase(context.Context, int64) error
	DeleteFromCart(context.Context, int64, uint32, uint32) error
}

type CartRepository interface {
	AddToCart(ctx context.Context, userID int64, sku uint32, count uint32) error
	UpdateCountCart(ctx context.Context, userID int64, sku uint32, count uint32) error
	DeleteCart(ctx context.Context, userID int64, sku uint32) error
	DeleteUserCart(ctx context.Context, userID int64) error
	ListCart(ctx context.Context, userID int64) ([]model.CartItem, error)
	GetCartItem(ctx context.Context, userID int64, sku uint32) (model.CartItem, error)
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
	// todo RunSerializable()
	// todo RunCommitted()
	// todo RunUncommitted()
}

var _ BusinessLogic = (*domain)(nil)

type domain struct {
	lomsService        LomsService
	productService     ProductService
	transactionManager TransactionManager
	cartRepository     CartRepository
}

func New(
	lomsService LomsService,
	productService ProductService,
	transactionManager TransactionManager,
	cartRepository CartRepository,
) *domain {
	return &domain{
		lomsService:        lomsService,
		productService:     productService,
		transactionManager: transactionManager,
		cartRepository:     cartRepository,
	}
}

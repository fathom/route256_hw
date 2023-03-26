package domain

import (
	"context"
	"route256/checkout/internal/clients/grpc/loms_client"
	"route256/checkout/internal/clients/grpc/product_client"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/db_repository"
	"route256/checkout/internal/repository/db_repository/transactor"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i Limiter -o ./mocks/ -s "_minimock.go"

type Limiter interface {
	Wait(ctx context.Context) (err error)
}

type BusinessLogic interface {
	AddToCart(context.Context, int64, uint32, uint32) error
	ListCart(context.Context, int64) ([]model.CartItem, error)
	Purchase(context.Context, int64) error
	DeleteFromCart(context.Context, int64, uint32, uint32) error
}

var _ BusinessLogic = (*domain)(nil)

type domain struct {
	lomsService        loms_client.LomsService
	productService     product_client.ProductService
	transactionManager transactor.TransactionManager
	cartRepository     db_repository.CartRepository
	limiter            Limiter
}

func New(
	lomsService loms_client.LomsService,
	productService product_client.ProductService,
	transactionManager transactor.TransactionManager,
	cartRepository db_repository.CartRepository,
	limiter Limiter,
) *domain {
	return &domain{
		lomsService:        lomsService,
		productService:     productService,
		transactionManager: transactionManager,
		cartRepository:     cartRepository,
		limiter:            limiter,
	}
}

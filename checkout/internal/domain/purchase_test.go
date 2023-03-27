package domain

import (
	"context"
	"errors"
	lc "route256/checkout/internal/clients/grpc/loms_client"
	lcMocks "route256/checkout/internal/clients/grpc/loms_client/mocks"
	pc "route256/checkout/internal/clients/grpc/product_client"
	pcMocks "route256/checkout/internal/clients/grpc/product_client/mocks"
	domainMocks "route256/checkout/internal/domain/mocks"
	"route256/checkout/internal/model"
	cr "route256/checkout/internal/repository/db_repository"
	crMocks "route256/checkout/internal/repository/db_repository/mocks"
	ts "route256/checkout/internal/repository/db_repository/transactor"
	tsMocks "route256/checkout/internal/repository/db_repository/transactor/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
)

func TestPurchase(t *testing.T) {

	type cartRepositoryMockFunc func(mc *minimock.Controller) cr.CartRepository
	type transactionManagerMockFunc func(mc *minimock.Controller) ts.TransactionManager
	type productServiceMockFunc func(mc *minimock.Controller) pc.ProductService
	type lomsServiceMockFunc func(mc *minimock.Controller) lc.LomsService
	type limiterMockFunc func(mc *minimock.Controller) Limiter

	type args struct {
		ctx    context.Context
		userID int64
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()
		n   = 5

		ListCartResponse []model.CartItem
		OrderItems       []*model.OrderItem

		userID  = gofakeit.Int64()
		orderID = gofakeit.Int64()

		clientErr = errors.New("client error")
	)

	for i := 0; i < n; i++ {
		ListCartResponse = append(ListCartResponse, model.CartItem{
			Sku:   gofakeit.Uint32(),
			Count: gofakeit.Uint32(),
			Name:  gofakeit.Name(),
			Price: gofakeit.Uint32(),
		})
	}

	for _, item := range ListCartResponse {
		OrderItems = append(OrderItems, &model.OrderItem{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	tests := []struct {
		name        string
		args        args
		err         error
		wantErr     bool
		crMock      cartRepositoryMockFunc
		trMock      transactionManagerMockFunc
		pcMock      productServiceMockFunc
		lcMock      lomsServiceMockFunc
		limiterMock limiterMockFunc
	}{
		{
			name: "Test successful",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			err:     nil,
			wantErr: false,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return(ListCartResponse, nil)
				mock.DeleteUserCartMock.Expect(ctx, userID).Return(nil)
				return mock
			},
			trMock: func(mc *minimock.Controller) ts.TransactionManager {
				mock := tsMocks.NewTransactionManagerMock(mc)
				return mock
			},
			pcMock: func(mc *minimock.Controller) pc.ProductService {
				mock := pcMocks.NewProductServiceMock(mc)
				return mock
			},
			lcMock: func(mc *minimock.Controller) lc.LomsService {
				mock := lcMocks.NewLomsServiceMock(mc)
				mock.CreateOrderMock.Expect(ctx, userID, OrderItems).Return(orderID, nil)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) Limiter {
				mock := domainMocks.NewLimiterMock(mc)
				return mock
			},
		},
		{
			name: "Test get error from cartRepository.ListCart",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			err:     cr.ErrCartRepository,
			wantErr: true,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return([]model.CartItem{}, cr.ErrCartRepository)
				mock.DeleteUserCartMock.Expect(ctx, userID).Return(nil)
				return mock
			},
			trMock: func(mc *minimock.Controller) ts.TransactionManager {
				mock := tsMocks.NewTransactionManagerMock(mc)
				return mock
			},
			pcMock: func(mc *minimock.Controller) pc.ProductService {
				mock := pcMocks.NewProductServiceMock(mc)
				return mock
			},
			lcMock: func(mc *minimock.Controller) lc.LomsService {
				mock := lcMocks.NewLomsServiceMock(mc)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) Limiter {
				mock := domainMocks.NewLimiterMock(mc)
				return mock
			},
		},
		{
			name: "Test get error from cartRepository.DeleteUserCart",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			err:     cr.ErrCartRepository,
			wantErr: true,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return(ListCartResponse, nil)
				mock.DeleteUserCartMock.Expect(ctx, userID).Return(cr.ErrCartRepository)
				return mock
			},
			trMock: func(mc *minimock.Controller) ts.TransactionManager {
				mock := tsMocks.NewTransactionManagerMock(mc)
				return mock
			},
			pcMock: func(mc *minimock.Controller) pc.ProductService {
				mock := pcMocks.NewProductServiceMock(mc)
				return mock
			},
			lcMock: func(mc *minimock.Controller) lc.LomsService {
				mock := lcMocks.NewLomsServiceMock(mc)
				mock.CreateOrderMock.Expect(ctx, userID, OrderItems).Return(orderID, nil)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) Limiter {
				mock := domainMocks.NewLimiterMock(mc)
				return mock
			},
		},
		{
			name: "Test get error from lomsService",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			err:     clientErr,
			wantErr: true,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return(ListCartResponse, nil)
				mock.DeleteUserCartMock.Expect(ctx, userID).Return(nil)
				return mock
			},
			trMock: func(mc *minimock.Controller) ts.TransactionManager {
				mock := tsMocks.NewTransactionManagerMock(mc)
				return mock
			},
			pcMock: func(mc *minimock.Controller) pc.ProductService {
				mock := pcMocks.NewProductServiceMock(mc)
				return mock
			},
			lcMock: func(mc *minimock.Controller) lc.LomsService {
				mock := lcMocks.NewLomsServiceMock(mc)
				mock.CreateOrderMock.Expect(ctx, userID, OrderItems).Return(0, clientErr)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) Limiter {
				mock := domainMocks.NewLimiterMock(mc)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &domain{
				lomsService:        tt.lcMock(mc),
				productService:     tt.pcMock(mc),
				transactionManager: tt.trMock(mc),
				cartRepository:     tt.crMock(mc),
				limiter:            tt.limiterMock(mc),
			}
			if err := d.Purchase(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("Purchase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

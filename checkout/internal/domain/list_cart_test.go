package domain

import (
	"context"
	"errors"
	"reflect"
	lc "route256/checkout/internal/clients/grpc/loms_client"
	lcMocks "route256/checkout/internal/clients/grpc/loms_client/mocks"
	pc "route256/checkout/internal/clients/grpc/product_client"
	pcMocks "route256/checkout/internal/clients/grpc/product_client/mocks"
	"route256/checkout/internal/config"
	domainMocks "route256/checkout/internal/domain/mocks"
	"route256/checkout/internal/model"
	cr "route256/checkout/internal/repository/db_repository"
	crMocks "route256/checkout/internal/repository/db_repository/mocks"
	ts "route256/checkout/internal/repository/db_repository/transactor"
	tsMocks "route256/checkout/internal/repository/db_repository/transactor/mocks"
	"route256/libs/workerpool"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
)

func TestListCart(t *testing.T) {
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

		n = 5

		ListCartResponse []model.CartItem

		userID = gofakeit.Int64()

		clientErr  = errors.New("client error")
		limiterErr = errors.New("limiter error")
	)
	defer mc.Finish()

	ctxWithCancel, cancel := context.WithDeadline(ctx, time.Now().Add(-7*time.Hour))
	cancel()

	for i := 0; i < n; i++ {
		ListCartResponse = append(ListCartResponse, model.CartItem{
			Sku:   gofakeit.Uint32(),
			Count: gofakeit.Uint32(),
			Name:  gofakeit.Name(),
			Price: gofakeit.Uint32(),
		})
	}

	tests := []struct {
		name         string
		args         args
		want         []model.CartItem
		err          error
		wantErr      bool
		crMock       cartRepositoryMockFunc
		trMock       transactionManagerMockFunc
		pcMock       productServiceMockFunc
		lcMock       lomsServiceMockFunc
		limiterMock  limiterMockFunc
		countWorkers int
	}{
		{
			name: "Test successful",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want:    ListCartResponse,
			err:     nil,
			wantErr: false,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return(ListCartResponse, nil)
				return mock
			},
			trMock: func(mc *minimock.Controller) ts.TransactionManager {
				mock := tsMocks.NewTransactionManagerMock(mc)
				return mock
			},
			pcMock: func(mc *minimock.Controller) pc.ProductService {
				mock := pcMocks.NewProductServiceMock(mc)
				for _, item := range ListCartResponse {
					mock.GetProductMock.When(ctx, item.Sku).Then(item.Name, item.Price, nil)
				}
				return mock
			},
			lcMock: func(mc *minimock.Controller) lc.LomsService {
				mock := lcMocks.NewLomsServiceMock(mc)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) Limiter {
				mock := domainMocks.NewLimiterMock(mc)
				mock.WaitMock.Expect(ctx).Return(nil)
				return mock
			},
			countWorkers: 1,
		},
		{
			name: "Test get error from cartRepository.ListCart",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want:    nil,
			err:     cr.ErrCartRepository,
			wantErr: true,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return([]model.CartItem{}, cr.ErrCartRepository)
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
			countWorkers: 1,
		},
		{
			name: "Test get error from ProductService.GetProduct",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want:    nil,
			err:     clientErr,
			wantErr: true,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return(ListCartResponse, nil)
				return mock
			},
			trMock: func(mc *minimock.Controller) ts.TransactionManager {
				mock := tsMocks.NewTransactionManagerMock(mc)
				return mock
			},
			pcMock: func(mc *minimock.Controller) pc.ProductService {
				mock := pcMocks.NewProductServiceMock(mc)
				mock.GetProductMock.Return("", 0, clientErr)
				return mock
			},
			lcMock: func(mc *minimock.Controller) lc.LomsService {
				mock := lcMocks.NewLomsServiceMock(mc)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) Limiter {
				mock := domainMocks.NewLimiterMock(mc)
				mock.WaitMock.Expect(ctx).Return(nil)
				return mock
			},
			countWorkers: 1,
		},
		{
			name: "Test get error from Limiter",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want:    nil,
			err:     limiterErr,
			wantErr: true,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, userID).Return(ListCartResponse, nil)
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
				mock.WaitMock.Return(limiterErr)
				return mock
			},
			countWorkers: 1,
		},
		{
			name: "Test get cancel Context",
			args: args{
				ctx:    ctxWithCancel,
				userID: userID,
			},
			want:    []model.CartItem{},
			err:     nil,
			wantErr: false,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Return(ListCartResponse, nil)
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
			countWorkers: 1,
		},
		{
			name: "Test with 0 worker config",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want:    nil,
			err:     workerpool.ErrNoWorkers,
			wantErr: true,
			crMock: func(mc *minimock.Controller) cr.CartRepository {
				mock := crMocks.NewCartRepositoryMock(mc)
				mock.ListCartMock.Return(ListCartResponse, nil)
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
			countWorkers: 0,
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
			config.ConfigData.CountWorkers = tt.countWorkers
			got, err := d.ListCart(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListCart() got = %v, want %v", got, tt.want)
			}
		})
	}
}

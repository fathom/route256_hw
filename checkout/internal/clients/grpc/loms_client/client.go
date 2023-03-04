package loms_client

import (
	"context"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/model"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc"
)

var _ domain.LomsService = &client{}

type client struct {
	lomsService desc.LomsV1Client
}

func New(cc *grpc.ClientConn) *client {
	return &client{
		desc.NewLomsV1Client(cc),
	}
}

func (c *client) Stocks(ctx context.Context, sku uint32) ([]*model.Stock, error) {
	stocksResponse, err := c.lomsService.Stocks(ctx, &desc.StocksRequest{Sku: sku})
	if err != nil {
		return nil, err
	}

	stocks := make([]*model.Stock, 0, len(stocksResponse.GetStocks()))
	for _, stock := range stocksResponse.GetStocks() {
		stocks = append(stocks, &model.Stock{
			WarehouseID: stock.GetWarehouseId(),
			Count:       stock.GetCount(),
		})
	}

	return stocks, nil
}

func (c *client) CreateOrder(ctx context.Context, user int64, items []*model.OrderItem) (int64, error) {
	createOrderResponse, err := c.lomsService.CreateOrder(ctx, &desc.CreateOrderRequest{
		User:  user,
		Items: converter.ToOrderItemListDesc(items),
	})
	if err != nil {
		return 0, err
	}

	return createOrderResponse.GetOrderId(), nil
}

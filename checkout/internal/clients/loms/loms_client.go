package loms

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/clientwrapper"
)

type Client struct {
	url            string
	urlStocks      string
	urlCreateOrder string
}

func New(url string) *Client {
	return &Client{
		url:            url,
		urlStocks:      url + "/stocks",
		urlCreateOrder: url + "/createOrder",
	}
}

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type StocksResponse struct {
	Stocks []StocksItem `json:"stocks"`
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	request := StocksRequest{SKU: sku}
	var response StocksResponse
	if err := clientwrapper.New(request, &response, c.urlStocks).DoRequest(ctx); err != nil {
		return nil, err
	}

	stocks := make([]domain.Stock, 0, len(response.Stocks))
	for _, stock := range response.Stocks {
		stocks = append(stocks, domain.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return stocks, nil
}

type CreateOrderRequest struct {
	User  int64              `json:"user"`
	Items []domain.OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderID int64 `json:"orderID"`
}

func (c *Client) CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
	request := CreateOrderRequest{
		User:  user,
		Items: items,
	}
	var response CreateOrderResponse

	if err := clientwrapper.New(request, &response, c.urlCreateOrder).DoRequest(ctx); err != nil {
		return 0, err
	}

	return response.OrderID, nil
}

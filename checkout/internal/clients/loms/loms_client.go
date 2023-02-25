package loms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"route256/checkout/internal/domain"

	"github.com/pkg/errors"
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

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlStocks, bytes.NewBuffer(rawJSON))
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response StocksResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json")
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

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return 0, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlCreateOrder, bytes.NewBuffer(rawJSON))
	if err != nil {
		return 0, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return 0, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response CreateOrderResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return 0, errors.Wrap(err, "decoding json")
	}

	return response.OrderID, nil
}

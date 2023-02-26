package stocks

import (
	"context"
	"errors"
	"log"
)

//stocks
//Возвращает количество товаров, которые можно купить с разных складов.
//Если товар был зарезерванован у кого-то в заказе и ждет оплаты,
//его купить нельзя.

type Handler struct{}

type Request struct {
	SKU uint32 `json:"sku"`
}

type Response struct {
	Stocks []Stock `json:"stocks"`
}

type Stock struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

var (
	ErrEmptySku = errors.New("empty sku")
)

func (r Request) Validate() error {
	if r.SKU == 0 {
		return ErrEmptySku
	}

	return nil
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(_ context.Context, request Request) (Response, error) {
	log.Printf("stocks: %+v", request)

	var stocks []Stock
	stocks = append(stocks, Stock{
		WarehouseID: 1,
		Count:       10,
	}, Stock{
		WarehouseID: 2,
		Count:       20,
	})

	response := Response{
		Stocks: stocks,
	}

	return response, nil
}

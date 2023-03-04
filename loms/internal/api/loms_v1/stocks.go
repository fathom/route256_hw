package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"
)

//stocks
//Возвращает количество товаров, которые можно купить с разных складов.
//Если товар был зарезерванован у кого-то в заказе и ждет оплаты,
//его купить нельзя.

func (h *Handlers) Stocks(_ context.Context, request *desc.StocksRequest) (*desc.StocksResponse, error) {
	log.Printf("stocks: %+v", request)

	var stocks []*desc.StockItem

	stocks = append(stocks, &desc.StockItem{
		WarehouseId: 1,
		Count:       100,
	}, &desc.StockItem{
		WarehouseId: 2,
		Count:       20,
	})

	return &desc.StocksResponse{
		Stocks: stocks,
	}, nil
}

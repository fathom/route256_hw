package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"
)

// stocks
// Возвращает количество товаров, которые можно купить с разных складов.
// Если товар был зарезерванован у кого-то в заказе и ждет оплаты,
// его купить нельзя.

func (h *Handlers) Stocks(ctx context.Context, request *desc.StocksRequest) (*desc.StocksResponse, error) {
	log.Printf("stocks: %+v", request)

	items, err := h.businessLogic.Stocks(ctx, request.GetSku())
	if err != nil {
		return nil, err
	}

	var stocks []*desc.StockItem

	for _, item := range items {
		stocks = append(stocks, &desc.StockItem{
			WarehouseId: item.WarehouseID,
			Count:       uint64(item.Count),
		})
	}

	return &desc.StocksResponse{
		Stocks: stocks,
	}, nil
}

package domain

import (
	"context"
	"route256/loms/internal/model"
)

func (d *Domain) Stocks(ctx context.Context, sku uint32) ([]model.StockItem, error) {

	var result []model.StockItem

	result = append(result, model.StockItem{
		WarehouseId: 1,
		Count:       50,
	}, model.StockItem{
		WarehouseId: 2,
		Count:       10,
	})

	return result, nil
}

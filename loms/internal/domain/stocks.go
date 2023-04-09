package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/logger"
	"route256/loms/internal/model"
)

func (d *Domain) Stocks(ctx context.Context, sku uint32) ([]model.StockItem, error) {

	var result []model.StockItem

	stocks, err := d.WarehouseRepository.GetStocksBySku(ctx, sku)
	if err != nil {
		return nil, err
	}

	for _, stock := range stocks {
		count := stock.Count - stock.Reservation
		if count > 0 {
			result = append(result, model.StockItem{
				WarehouseID: stock.WarehouseID,
				Count:       count,
			})
		}
	}

	logger.Debug(fmt.Sprintf("Stocks available for %v", sku))

	return result, nil
}

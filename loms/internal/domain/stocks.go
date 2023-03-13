package domain

import (
	"context"
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

	return result, nil
}

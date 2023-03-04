package converter

import (
	"route256/checkout/internal/model"
	desc "route256/loms/pkg/loms_v1"
)

func ToOrderItemDesc(item *model.OrderItem) *desc.OrderItem {
	return &desc.OrderItem{
		Sku:   item.Sku,
		Count: item.Count,
	}
}

func ToOrderItemListDesc(items []*model.OrderItem) []*desc.OrderItem {
	res := make([]*desc.OrderItem, 0, len(items))
	for _, item := range items {
		res = append(res, ToOrderItemDesc(item))
	}
	return res
}

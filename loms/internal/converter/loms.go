package converter

import (
	"route256/loms/internal/model"
	desc "route256/loms/pkg/loms_v1"
)

func ToOrderItemDomain(item *desc.OrderItem) *model.OrderItem {
	return &model.OrderItem{
		Sku:   item.GetSku(),
		Count: item.GetCount(),
	}
}

func ToOrderItemListDomain(items []*desc.OrderItem) []*model.OrderItem {
	res := make([]*model.OrderItem, 0, len(items))
	for _, item := range items {
		res = append(res, ToOrderItemDomain(item))
	}
	return res
}

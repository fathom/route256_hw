package converter

import (
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

func ToOrderItemDomain(item *desc.OrderItem) *domain.OrderItem {
	return &domain.OrderItem{
		Sku:   item.GetSku(),
		Count: item.GetCount(),
	}
}

func ToOrderItemListDomain(items []*desc.OrderItem) []*domain.OrderItem {
	res := make([]*domain.OrderItem, 0, len(items))
	for _, item := range items {
		res = append(res, ToOrderItemDomain(item))
	}
	return res
}

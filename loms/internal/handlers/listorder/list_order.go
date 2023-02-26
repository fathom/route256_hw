package listorder

import (
	"context"
	"log"
	"route256/loms/internal/domain"
)

//listOrder
//Показывает информацию по заказу

type Handler struct{}

type Request struct {
	OrderID int64 `json:"orderID"`
}

type Response struct {
	Status domain.OrderStatus `json:"status"`
	User   int64              `json:"user"`
	Items  []domain.OrderItem `json:"items"`
}

func (r Request) Validate() error {
	return nil
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(_ context.Context, request Request) (Response, error) {
	log.Printf("listOrder: %+v", request)

	var items []domain.OrderItem
	items = append(items, domain.OrderItem{
		Sku:   1,
		Count: 10,
	}, domain.OrderItem{
		Sku:   2,
		Count: 20,
	})

	response := Response{
		Status: domain.AwaitingPayment,
		User:   100,
		Items:  items,
	}

	_ = domain.NewStatus
	_ = domain.Cancelled
	_ = domain.Failed
	_ = domain.Payed

	return response, nil
}

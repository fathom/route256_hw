package cancelorder

import (
	"context"
	"errors"
	"log"
)

//cancelOrder
//Отменяет заказ, снимает резерв со всех товаров в заказе.

type Handler struct{}

type Request struct {
	OrderID int64 `json:"orderID"`
}

type Response struct{}

var (
	ErrEmptyOrderID = errors.New("empty order id")
)

func (r Request) Validate() error {
	if r.OrderID == 0 {
		return ErrEmptyOrderID
	}

	return nil
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(_ context.Context, request Request) (Response, error) {
	log.Printf("cancelOrder: %+v", request)

	var response Response

	return response, nil
}

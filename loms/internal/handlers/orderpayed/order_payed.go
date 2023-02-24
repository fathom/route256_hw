package orderpayed

import (
	"context"
	"errors"
	"log"
)

//orderPayed
//Помечает заказ оплаченным. Зарезервированные товары должны перейти
//в статус купленных.

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
	log.Printf("orderPayed: %+v", request)

	var response Response

	return response, nil
}

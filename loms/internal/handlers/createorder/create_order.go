package createorder

import (
	"context"
	"log"
	"route256/loms/internal/domain"
)

//createOrder
//Создает новый заказ для пользователя из списка переданных товаров.
//Товары при этом нужно зарезервировать на складе.

type Handler struct {
	businessLogic *domain.Domain
}

type Request struct {
	User  int64              `json:"user"`
	Items []domain.OrderItem `json:"items"`
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

func (r Request) Validate() error {
	return nil
}

func New(businessLogic *domain.Domain) *Handler {
	return &Handler{
		businessLogic: businessLogic,
	}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("createOrder: %+v", request)

	var response Response
	orderID, err := h.businessLogic.CreateOrder(ctx, request.User, request.Items)
	if err != nil {
		return response, err
	}

	response = Response{
		OrderID: orderID,
	}

	return response, nil
}

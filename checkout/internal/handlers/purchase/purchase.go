package purchase

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
)

//purchase
//Оформить заказ по всем товарам корзины. Вызывает createOrder у LOMS.

type Handler struct {
	businessLogic *domain.Domain
}

func New(businessLogic *domain.Domain) *Handler {
	return &Handler{
		businessLogic: businessLogic,
	}
}

type Request struct {
	User int64 `json:"user"`
}

type Response struct{}

var (
	ErrEmptyUser = errors.New("empty user")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("purchase: %+v", request)

	var response Response

	err := h.businessLogic.Purchase(ctx, request.User)
	if err != nil {
		return response, err
	}

	return response, nil
}

package addtocart

import (
	"context"
	"errors"
	"log"
)

//addToCart
//Добавить товар в корзину определенного пользователя. При этом надо
//проверить наличие товара через LOMS.stocks

type BusinessLogic interface {
	AddToCart(context.Context, int64, uint32, uint16) error
}

type Handler struct {
	businessLogic BusinessLogic
}

func New(businessLogic BusinessLogic) *Handler {
	return &Handler{
		businessLogic: businessLogic,
	}
}

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Response struct{}

var (
	ErrEmptyUser = errors.New("empty user")
	ErrEmptySKU  = errors.New("empty sku")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	if r.SKU == 0 {
		return ErrEmptySKU
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("addToCart: %+v", request)

	var response Response

	err := h.businessLogic.AddToCart(ctx, request.User, request.SKU, request.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

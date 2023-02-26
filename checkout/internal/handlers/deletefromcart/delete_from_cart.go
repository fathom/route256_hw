package deletefromcart

import (
	"context"
	"errors"
	"log"
)

//deleteFromCart
//Удалить товар из корзины определенного пользователя.

type Handler struct{}

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

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(_ context.Context, request Request) (Response, error) {
	log.Printf("deleteFromCart: %+v", request)

	var response Response

	return response, nil
}

package listcart

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
)

//listCart
//Показать список товаров в корзине с именами и ценами (их надо в реальном
//времени получать из ProductService)

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

type Product struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Response struct {
	Items      []Product `json:"items"`
	TotalPrice uint32    `json:"totalPrice"`
}

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
	log.Printf("listCart: %+v", request)

	var response Response

	cartItems, err := h.businessLogic.ListCart(ctx, request.User)
	if err != nil {
		return response, err
	}

	for _, item := range cartItems {
		response.Items = append(response.Items, Product{item.Scu, item.Count, item.Name, item.Price})
		response.TotalPrice += item.Price
	}

	return response, nil
}

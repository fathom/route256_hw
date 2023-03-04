package checkout_v1

import (
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
)

type Handlers struct {
	desc.UnimplementedCheckoutV1Server
	businessLogic domain.BusinessLogic
}

func NewCheckoutV1(businessLogic domain.BusinessLogic) *Handlers {
	return &Handlers{
		businessLogic: businessLogic,
	}
}

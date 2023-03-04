package loms_client

import (
	"context"
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/product"

	"google.golang.org/grpc"
)

var _ domain.ProductService = &client{}

type client struct {
	productService desc.ProductServiceClient
	token          string
}

func New(cc *grpc.ClientConn, token string) *client {
	return &client{
		desc.NewProductServiceClient(cc),
		token,
	}
}

func (c *client) GetProduct(ctx context.Context, sku uint32) (string, uint32, error) {
	response, err := c.productService.GetProduct(ctx, &desc.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	})
	if err != nil {
		return "", 0, err
	}

	return response.GetName(), response.GetPrice(), nil
}

func (c *client) ListSkus(ctx context.Context, startAfterSku, count uint32) ([]uint32, error) {
	response, err := c.productService.ListSkus(ctx, &desc.ListSkusRequest{
		Token:         c.token,
		StartAfterSku: startAfterSku,
		Count:         count,
	})
	if err != nil {
		return nil, err
	}

	return response.GetSkus(), nil
}

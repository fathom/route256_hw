package product_client

import (
	"context"
	"go.uber.org/zap"
	c "route256/checkout/internal/cache"
	"route256/checkout/internal/logger"
	"route256/checkout/internal/metrics"
	desc "route256/product_service/pkg/product_service"
	"time"

	"google.golang.org/grpc"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ProductService -o ./mocks/ -s "_minimock.go"

type ProductService interface {
	GetProduct(ctx context.Context, sku uint32) (string, uint32, error)
	ListSkus(ctx context.Context, startAfterSku, count uint32) ([]uint32, error)
}

var _ ProductService = &client{}

type client struct {
	productService desc.ProductServiceClient
	cache          c.Cache[uint32, *desc.GetProductResponse]
	token          string
}

func New(cc *grpc.ClientConn, cache c.Cache[uint32, *desc.GetProductResponse], token string) *client {
	return &client{
		desc.NewProductServiceClient(cc),
		cache,
		token,
	}
}

func (c *client) GetProduct(ctx context.Context, sku uint32) (string, uint32, error) {

	response, ok := c.cache.Get(sku)
	metrics.CacheRequestCounter.Inc()
	logger.Debug("Request cached sku", zap.Uint32("sku", sku))

	if !ok {
		metrics.CacheErrorsCounter.Inc()
		logger.Debug("Miss cached sku", zap.Uint32("sku", sku))
		var err error
		response, err = c.productService.GetProduct(ctx, &desc.GetProductRequest{
			Token: c.token,
			Sku:   sku,
		})
		if err != nil {
			return "", 0, err
		}

		c.cache.Set(sku, response, time.Hour)
	} else {
		metrics.CacheHitsCounter.Inc()
		logger.Debug("Hit cached sku", zap.Uint32("sku", sku))
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

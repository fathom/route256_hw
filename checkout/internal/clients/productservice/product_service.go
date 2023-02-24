package productservice

import (
	"context"
	"route256/libs/clientwrapper"
)

//Swagger развернут по адресу: http://route256.pavl.uk:8080/docs/
//GRPC развернуто по адресу: route256.pavl.uk:8082
//get_product
//list_skus

type Client struct {
	url           string
	urlGetProduct string
	urlListSkus   string
	token         string
}

func New(url, token string) *Client {
	return &Client{
		url:           url,
		urlGetProduct: url + "/get_product",
		urlListSkus:   url + "/list_skus",
		token:         token,
	}
}

type ProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type ProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, sku uint32) (string, uint32, error) {
	request := ProductRequest{
		Token: c.token,
		SKU:   sku,
	}
	var response ProductResponse
	if err := clientwrapper.New(request, &response, c.urlGetProduct).DoRequest(ctx); err != nil {
		return "", 0, err
	}

	return response.Name, response.Price, nil
}

type ListRequest struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}

type ListResponse struct {
	Skus []uint32 `json:"skus"`
}

func (c *Client) ListSkus(ctx context.Context, startAfterSku, count uint32) ([]uint32, error) {
	request := ListRequest{
		Token:         c.token,
		StartAfterSku: startAfterSku,
		Count:         count,
	}
	var response ListResponse
	if err := clientwrapper.New(request, &response, c.urlListSkus).DoRequest(ctx); err != nil {
		return nil, err
	}

	return response.Skus, nil
}

package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
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

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return "", 0, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlGetProduct, bytes.NewBuffer(rawJSON))
	if err != nil {
		return "", 0, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return "", 0, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var serviceResponse ProductResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&serviceResponse)
	if err != nil {
		return "", 0, errors.Wrap(err, "decoding json")
	}

	return serviceResponse.Name, serviceResponse.Price, nil
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

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlListSkus, bytes.NewBuffer(rawJSON))
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var serviceResponse ListResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&serviceResponse)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json")
	}

	return serviceResponse.Skus, nil
}

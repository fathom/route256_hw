package clientwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Wrapper[Request any, Response any] struct {
	request  Request
	response *Response
	url      string
}

func New[Request any, Response any](request Request, response *Response, url string) *Wrapper[Request, Response] {
	return &Wrapper[Request, Response]{
		request:  request,
		response: response,
		url:      url,
	}
}

func (wrapper *Wrapper[Request, Response]) DoRequest(ctx context.Context) error {
	rawJSON, err := json.Marshal(&wrapper.request)
	if err != nil {
		return errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, wrapper.url, bytes.NewBuffer(rawJSON))
	if err != nil {
		return errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&wrapper.response)
	if err != nil {
		return errors.Wrap(err, "decoding json")
	}

	return nil
}

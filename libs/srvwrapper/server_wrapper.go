package srvwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Validator interface {
	Validate() error
}

type Wrapper[Request Validator, Response any] struct {
	fn func(ctx context.Context, request Request) (Response, error)
}

func New[Request Validator, Response any](fn func(ctx context.Context, request Request) (Response, error)) *Wrapper[Request, Response] {
	return &Wrapper[Request, Response]{
		fn: fn,
	}
}

func (wrapper *Wrapper[Request, Response]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limitedReader := io.LimitReader(r.Body, 1_000_000)

	var request Request
	err := json.NewDecoder(limitedReader).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorText(w, "decoding JSON", err)
		return
	}

	err = request.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorText(w, "validating request", err)
		return
	}

	response, err := wrapper.fn(ctx, request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorText(w, "running handler", err)
		return
	}

	rawJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorText(w, "encoding JSON", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(rawJSON)
}

func writeErrorText(w http.ResponseWriter, text string, err error) {
	buf := bytes.NewBufferString(text)

	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteByte('\n')

	_, _ = w.Write(buf.Bytes())
}

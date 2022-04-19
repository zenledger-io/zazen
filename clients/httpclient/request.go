package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
)

type Request struct {
	URL     string
	Method  string
	Body    io.Reader
	Headers map[string]string
	Query   map[string]string
	Opts    []RequestOption
}

func (r *Request) Options() []RequestOption {
	opts := make([]RequestOption, 0, len(r.Opts)+2)

	if len(r.Opts) > 0 {
		opts = append(opts, r.Opts...)
	}

	if len(r.Query) > 0 {
		opts = append(opts, RequestSetQuery{Query: r.Query})
	}

	if len(r.Headers) > 0 {
		opts = append(opts, RequestSetHeaders{Headers: r.Headers})
	}

	return opts
}

func (r *Request) SetJSONBody(body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	r.Body = bytes.NewBuffer(b)
	r.Opts = append(r.Opts, RequestContentTypeJSON{})

	return nil
}

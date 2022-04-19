package httpclient

import "net/http"

var (
	ContentTypeJSON = "application/json; charset=utf-8"
)

type RequestOption interface {
	Configure(*http.Request) *http.Request
}

// RequestSetQuery calls Set on the request URL's query with each key value pair.
type RequestSetQuery struct {
	Query map[string]string
}

func (opt RequestSetQuery) Configure(r *http.Request) *http.Request {
	q := r.URL.Query()
	handleKeyValuesFunc(q.Set, opt.Query)
	r.URL.RawQuery = q.Encode()

	return r
}

// RequestAddQuery calls Add on the request URL's query with each key value pair.
type RequestAddQuery struct {
	Query map[string]string
}

func (opt RequestAddQuery) Configure(r *http.Request) *http.Request {
	q := r.URL.Query()
	handleKeyValuesFunc(q.Add, opt.Query)
	r.URL.RawQuery = q.Encode()

	return r
}

// RequestSetHeaders calls Set on the request URL's header with each key value pair.
type RequestSetHeaders struct {
	Headers map[string]string
}

func (opt RequestSetHeaders) Configure(r *http.Request) *http.Request {
	handleKeyValuesFunc(r.Header.Set, opt.Headers)

	return r
}

// RequestAddHeaders calls Add on the request URL's header with each key value pair.
type RequestAddHeaders struct {
	Headers map[string]string
}

func (opt RequestAddHeaders) Configure(r *http.Request) *http.Request {
	handleKeyValuesFunc(r.Header.Add, opt.Headers)

	return r
}

// RequestContentType sets the URL's Content-Type header to the supplied string.
type RequestContentType struct {
	ContentType string
}

func (opt RequestContentType) Configure(r *http.Request) *http.Request {
	r.Header.Set("Content-Type", opt.ContentType)
	return RequestSetHeaders{Headers: map[string]string{"Content-Type": opt.ContentType}}.Configure(r)
}

// RequestContentTypeJSON sets the URL's Content-Type header to "application/json; charset=utf-8".
type RequestContentTypeJSON struct{}

func (opt RequestContentTypeJSON) Configure(r *http.Request) *http.Request {
	return RequestContentType{ContentType: ContentTypeJSON}.Configure(r)
}

// helpers

func handleKeyValuesFunc(f func(string, string), m map[string]string) {
	for k, v := range m {
		f(k, v)
	}
}

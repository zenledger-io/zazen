package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewDefault()

	type Payload struct {
		Int int `json:"int"`
	}

	tcs := map[string]struct {
		Request Request
		Body    interface{}
		Opts    []RequestOption
		Handler http.HandlerFunc
	}{
		"GET request no overrides": {
			Request: Request{
				Method:  http.MethodGet,
				Headers: nil,
				Query:   nil,
				Opts: []RequestOption{
					RequestContentTypeJSON{},
					RequestSetQuery{Query: map[string]string{"q1": "a"}},
					RequestSetHeaders{Headers: map[string]string{"h1": "b"}},
				},
			},
			Body: nil,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "a" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "b" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
		"GET request override header and query request opts by setting them in request": {
			Request: Request{
				Method:  http.MethodGet,
				Headers: map[string]string{"h1": "c"},
				Query:   map[string]string{"q1": "d"},
				Opts: []RequestOption{
					RequestContentTypeJSON{},
					RequestSetQuery{Query: map[string]string{"q1": "a"}},
					RequestSetHeaders{Headers: map[string]string{"h1": "b"}},
				},
			},
			Body: nil,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "d" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "c" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
		"GET request override header and query everywhere by setting them in opts passed to SendRequest": {
			Request: Request{
				Method:  http.MethodGet,
				Headers: map[string]string{"h1": "c"},
				Query:   map[string]string{"q1": "d"},
				Opts: []RequestOption{
					RequestContentTypeJSON{},
					RequestSetQuery{Query: map[string]string{"q1": "a"}},
					RequestSetHeaders{Headers: map[string]string{"h1": "b"}},
				},
			},
			Body: nil,
			Opts: []RequestOption{
				RequestContentType{ContentType: "whatever"},
				RequestSetQuery{Query: map[string]string{"q1": "e"}},
				RequestSetHeaders{Headers: map[string]string{"h1": "f"}},
			},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "e" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "f" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != "whatever" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
		"POST request overrides request options content type when sending json": {
			Request: Request{
				Method: http.MethodPost,
				Opts: []RequestOption{
					RequestContentType{ContentType: "whatever"},
				},
			},
			Body: Payload{Int: 1},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				var p Payload
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if p != (Payload{Int: 1}) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			srv := httptest.NewServer(tc.Handler)
			defer srv.Close()

			tc.Request.URL = srv.URL
			require.NoError(t, tc.Request.SetJSONBody(tc.Body))

			resp, err := c.SendRequest(ctx, &tc.Request, tc.Opts...)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})
	}
}

func TestClient_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewDefault()

	tcs := map[string]struct {
		Query   map[string]string
		Headers map[string]string
		Opts    []RequestOption
		Handler http.HandlerFunc
	}{
		"when all request params are present": {
			Query:   map[string]string{"q1": "a"},
			Headers: map[string]string{"h1": "b"},
			Opts:    []RequestOption{RequestContentTypeJSON{}},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "a" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "b" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			srv := httptest.NewServer(tc.Handler)
			defer srv.Close()

			opts := append(tc.Opts, RequestSetQuery{Query: tc.Query}, RequestSetHeaders{Headers: tc.Headers})
			resp, err := c.Get(ctx, srv.URL, opts...)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})
	}
}

func TestClient_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewDefault()

	tcs := map[string]struct {
		Query   map[string]string
		Headers map[string]string
		Opts    []RequestOption
		Handler http.HandlerFunc
	}{
		"when all request params are present": {
			Query:   map[string]string{"q1": "a"},
			Headers: map[string]string{"h1": "b"},
			Opts:    []RequestOption{RequestContentTypeJSON{}},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "a" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "b" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			srv := httptest.NewServer(tc.Handler)
			defer srv.Close()

			opts := append(tc.Opts, RequestSetQuery{Query: tc.Query}, RequestSetHeaders{Headers: tc.Headers})
			resp, err := c.Delete(ctx, srv.URL, opts...)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})
	}
}

func TestClient_Post(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewDefault()

	type Payload struct {
		Int int `json:"int"`
	}

	tcs := map[string]struct {
		Query   map[string]string
		Headers map[string]string
		Payload Payload
		Opts    []RequestOption
		Handler http.HandlerFunc
	}{
		"when all request params are present": {
			Query:   map[string]string{"q1": "a"},
			Headers: map[string]string{"h1": "b"},
			Payload: Payload{Int: 1},
			Opts:    []RequestOption{RequestContentTypeJSON{}},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				var p Payload
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if p != (Payload{Int: 1}) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "a" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "b" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			srv := httptest.NewServer(tc.Handler)
			defer srv.Close()

			b, err := json.Marshal(tc.Payload)
			require.NoError(t, err)

			opts := append(tc.Opts, RequestSetQuery{Query: tc.Query}, RequestSetHeaders{Headers: tc.Headers})
			resp, err := c.Post(ctx, srv.URL, bytes.NewBuffer(b), opts...)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})
	}
}

func TestClient_Put(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewDefault()

	type Payload struct {
		Int int `json:"int"`
	}

	tcs := map[string]struct {
		Query   map[string]string
		Headers map[string]string
		Payload Payload
		Opts    []RequestOption
		Handler http.HandlerFunc
	}{
		"when all request params are present": {
			Query:   map[string]string{"q1": "a"},
			Headers: map[string]string{"h1": "b"},
			Payload: Payload{Int: 1},
			Opts:    []RequestOption{RequestContentTypeJSON{}},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				var p Payload
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if p != (Payload{Int: 1}) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "a" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "b" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			srv := httptest.NewServer(tc.Handler)
			defer srv.Close()

			b, err := json.Marshal(tc.Payload)
			require.NoError(t, err)

			opts := append(tc.Opts, RequestSetQuery{Query: tc.Query}, RequestSetHeaders{Headers: tc.Headers})
			resp, err := c.Put(ctx, srv.URL, bytes.NewBuffer(b), opts...)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})
	}
}

func TestClient_Patch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewDefault()

	type Payload struct {
		Int int `json:"int"`
	}

	tcs := map[string]struct {
		Query   map[string]string
		Headers map[string]string
		Payload Payload
		Opts    []RequestOption
		Handler http.HandlerFunc
	}{
		"when all request params are present": {
			Query:   map[string]string{"q1": "a"},
			Headers: map[string]string{"h1": "b"},
			Payload: Payload{Int: 1},
			Opts:    []RequestOption{RequestContentTypeJSON{}},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPatch {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				var p Payload
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if p != (Payload{Int: 1}) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query()
				if q.Get("q1") != "a" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("h1") != "b" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if r.Header.Get("Content-Type") != ContentTypeJSON {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			srv := httptest.NewServer(tc.Handler)
			defer srv.Close()

			b, err := json.Marshal(tc.Payload)
			require.NoError(t, err)

			opts := append(tc.Opts, RequestSetQuery{Query: tc.Query}, RequestSetHeaders{Headers: tc.Headers})
			resp, err := c.Patch(ctx, srv.URL, bytes.NewBuffer(b), opts...)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})
	}
}

package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	cases := []struct {
		Desc       string
		Error      Error
		StatusCode int
		Resp       errorResponse
	}{
		{
			Desc:       "empty values",
			Error:      Error{},
			StatusCode: http.StatusInternalServerError,
			Resp: errorResponse{
				Message: http.StatusText(http.StatusInternalServerError),
			},
		},
		{
			Desc: "empty values with status code",
			Error: Error{
				StatusCode: http.StatusTeapot,
			},
			StatusCode: http.StatusTeapot,
			Resp: errorResponse{
				Message: http.StatusText(http.StatusTeapot),
			},
		},
		{
			Desc: "only err",
			Error: Error{
				Err:        errors.New("err 111"),
				StatusCode: http.StatusTeapot,
			},
			StatusCode: http.StatusTeapot,
			Resp: errorResponse{
				Message: "err 111",
			},
		},
		{
			Desc: "only message",
			Error: Error{
				Message:    "msg 111",
				StatusCode: http.StatusTeapot,
			},
			StatusCode: http.StatusTeapot,
			Resp: errorResponse{
				Message: "msg 111",
			},
		},
		{
			Desc: "err and message",
			Error: Error{
				Message:    "msg 111",
				Err:        errors.New("err 111"),
				StatusCode: http.StatusTeapot,
			},
			StatusCode: http.StatusTeapot,
			Resp: errorResponse{
				Message: "err 111: msg 111",
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.Desc, func(t *testing.T) {
			w := httptest.NewRecorder()

			WriteJSON(context.TODO(), w, c.Error)

			res := w.Result()
			defer res.Body.Close()

			var resp errorResponse
			err := json.NewDecoder(res.Body).Decode(&resp)
			assert.NoError(t, err)

			assert.Equal(t, jsonContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, c.StatusCode, res.StatusCode)
			assert.Equal(t, c.Resp, resp)
		})
	}
}

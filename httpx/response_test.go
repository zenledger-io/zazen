package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	cases := []struct {
		Desc       string
		Response   Response
		StatusCode int
	}{
		{
			Desc:       "empty payload",
			Response:   Response{},
			StatusCode: http.StatusOK,
		},
		{
			Desc:       "non-empty payload",
			Response:   Response{Payload: "111"},
			StatusCode: http.StatusOK,
		},
		{
			Desc: "non-empty payload and status code",
			Response: Response{
				Payload:    "111",
				StatusCode: http.StatusTeapot,
			},
			StatusCode: http.StatusTeapot,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.Desc, func(t *testing.T) {
			w := httptest.NewRecorder()

			WriteJSON(context.TODO(), w, c.Response)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, jsonContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, c.StatusCode, res.StatusCode)
		})
	}
}

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_postHandler(t *testing.T) {
	storeURL["https://yoga.org/"] = "875910c4"

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		want   want
		body   string
		method string
	}{
		{
			name: "success",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
			body:   "https://yoga.org/",
			method: http.MethodPost,
		},
		{
			name: "incorrect method",
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   "https://yoga.org/",
			method: http.MethodGet,
		},
		{
			name: "request body is empty",
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   "",
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			rec := httptest.NewRecorder()
			postHandler(storeURL)(rec, req)

			res := rec.Result()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.NotNil(t, string(resBody))
		})
	}
}

func Test_getHandler(t *testing.T) {
	id := "875910c4"
	storeURL["https://yoga.org/"] = "875910c4"

	type want struct {
		code        int
		response    string
		contentType string
		header      string
	}

	tests := []struct {
		name    string
		want    want
		method  string
		request string
	}{
		{
			name: "success",
			want: want{
				code:   http.StatusTemporaryRedirect,
				header: storeURL[id],
			},
			method:  http.MethodGet,
			request: "/875910c4",
		},
		{
			name: "incorrect method",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				header:      "",
			},
			method:  http.MethodPost,
			request: "/875910c4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.request, nil)
			rec := httptest.NewRecorder()
			getHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.header, res.Header.Get("Location"))
		})
	}
}

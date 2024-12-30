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

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func Test_postHandler(t *testing.T) {
	storeURL["https://yoga.org/"] = "875910c4"

	ts := httptest.NewServer(ShortenerRoutes())
	defer ts.Close()

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
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
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
			resp, _ := testRequest(t, ts, test.method, "/", strings.NewReader(test.body))
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func Test_getHandler(t *testing.T) {
	id := "875910c4"
	storeURL["https://yoga.org/"] = "875910c4"

	ts := httptest.NewServer(ShortenerRoutes())
	defer ts.Close()

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
				code:        http.StatusMethodNotAllowed,
				contentType: "",
				header:      "",
			},
			method:  http.MethodPost,
			request: "/875910c4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, "/{id}", nil)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, test.want.header, resp.Header.Get("Location"))
		})
	}
}

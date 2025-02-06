package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nastez/shortener/internal/app/handlers/urlhandlers"
)

type MemoryStorage map[string]string

var storeURL = MemoryStorage{}

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	req.RequestURI = ""
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "")
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

	routes, err := ShortenerRoutes("", "", nil)
	if err != nil {
		return
	}

	ts := httptest.NewServer(routes)
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
				contentType: "text/plain; charset=utf-8",
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

	routes, err := ShortenerRoutes("", "", nil)
	if err != nil {
		return
	}

	ts := httptest.NewServer(routes)
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

func Test_shortenerHandler(t *testing.T) {
	storeURL["https://yoga.org/"] = "875910c4"

	routes, err := ShortenerRoutes("", "", nil)
	if err != nil {
		return
	}

	ts := httptest.NewServer(routes)
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
				contentType: "application/json",
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
				contentType: "application/json",
			},
			body:   "",
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, "/api/shorten", strings.NewReader(test.body))
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func Test_getPing(t *testing.T) {
	psDefault := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `shortener`, `pupupu`, `shortener`)

	routes, err := ShortenerRoutes("", psDefault, nil)
	if err != nil {
		return
	}

	ts := httptest.NewServer(routes)
	defer ts.Close()

	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name   string
		want   want
		method string
	}{
		{
			name: "success",
			want: want{
				code: http.StatusOK,
			},
			method: http.MethodGet,
		},
		{
			name: "incorrect method",
			want: want{
				code:     http.StatusMethodNotAllowed,
				response: "",
			},
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, "/ping", nil)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
		})
	}
}

func TestGzipCompression(t *testing.T) {
	//var storeURL = storeURL

	handler, err := urlhandlers.New(nil, nil, "http://localhost:8080", "")
	if err != nil {
		return
	}

	ts := httptest.NewServer(GzipMiddleware(handler.ShortenerHandler()))
	defer ts.Close()

	requestBody := `{"url":"https://yoga.org/"}`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", ts.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept-Encoding", "")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", ts.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		_, err = io.ReadAll(zr)
		require.NoError(t, err)
	})
}

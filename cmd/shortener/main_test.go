package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nastez/shortener/internal/storage"
	storeMock "github.com/Nastez/shortener/internal/store/mocks"
)

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
	ctrl := gomock.NewController(t)
	s := storeMock.NewMockStore(ctrl)

	//установим условие: при любом вызове метода Save не возвращались ошибки
	s.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	appInstance, err := newApp(s, "http://localhost:0007", "")
	if err != nil {
		assert.Error(t, err)
	}

	handler := appInstance.PostHandler()
	ts := httptest.NewServer(handler)
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
			name: "incorrect method GET",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   "https://yoga.org/",
			method: http.MethodGet,
		},
		{
			name: "incorrect method DELETE",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   "https://yoga.org/",
			method: http.MethodGet,
		},
		{
			name: "incorrect method PUT",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   "https://yoga.org/",
			method: http.MethodGet,
		},
		{
			name: "empty body",
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
	store := storage.MemoryStorage{}
	store["https://yoga.org/"] = "875910c4"

	ctrl := gomock.NewController(t)
	s := storeMock.NewMockStore(ctrl)

	//установим условие: при любом вызове метода Get не возвращались ошибки
	s.EXPECT().
		Get(gomock.Any(), id).
		Return("875910c4", false, nil).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	appInstance, err := newApp(store, "http://localhost:0007", "")
	if err != nil {
		assert.Error(t, err)
	}

	handler := appInstance.GetHandler()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
		header      string
	}

	tests := []struct {
		name   string
		want   want
		method string
	}{
		{
			name: "success",
			want: want{
				code:   http.StatusTemporaryRedirect,
				header: storage.MemoryStorage{}[id],
			},
			method: http.MethodGet,
		},
		{
			name: "incorrect method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				header:      "",
			},
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "/", nil)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", "111")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			// Создаем `ResponseRecorder`, чтобы записать ответ
			w := httptest.NewRecorder()

			// Вызываем обработчик
			handler(w, req)

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, test.want.code, w.Code, "Response code didn't match expected")
			assert.Equal(t, test.want.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, test.want.header, w.Header().Get("Location"))
		})
	}
}

func Test_shortenerHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := storeMock.NewMockStore(ctrl)

	//установим условие: при любом вызове метода Save не возвращались ошибки
	s.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	appInstance, err := newApp(s, "http://localhost:0007", "")
	if err != nil {
		assert.Error(t, err)
	}

	handler := appInstance.ShortenerHandler()
	ts := httptest.NewServer(handler)
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
			body:   `{"url":"https://yoga.org/"}`,
			method: http.MethodPost,
		},
		{
			name: "incorrect method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   `{"url":"https://yoga.org/"}`,
			method: http.MethodGet,
		},
		{
			name: "request body is empty",
			want: want{
				code:        http.StatusInternalServerError,
				response:    "",
				contentType: "",
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

	appInstance, err := newApp(nil, "http://localhost:0007", psDefault)
	if err != nil {
		assert.Error(t, err)
	}

	handler := appInstance.GetPing()
	ts := httptest.NewServer(handler)
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

func Test_postBatchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := storeMock.NewMockStore(ctrl)

	//установим условие: при любом вызове метода SaveBatch не возвращались ошибки
	s.EXPECT().
		SaveBatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	appInstance, err := newApp(s, "http://localhost:0007", "")
	if err != nil {
		assert.Error(t, err)
	}

	handler := appInstance.PostBatch()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	body := `[
    {
        "correlation_id": "1",
        "original_url": "http://to1ghmjtw0f.biz"
    },
    {
        "correlation_id": "2",
        "original_url": "http://to2ghmjtw0f.biz"
    }
]`

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
			body:   body,
			method: http.MethodPost,
		},
		{
			name: "incorrect method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   body,
			method: http.MethodGet,
		},
		{
			name: "request body is empty",
			want: want{
				code:        http.StatusInternalServerError,
				response:    "",
				contentType: "",
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

func Test_getURLSHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := storeMock.NewMockStore(ctrl)

	//установим условие: при любом вызове метода GetURLs не возвращались ошибки
	s.EXPECT().
		GetURLs(gomock.Any(), gomock.Any()).
		Return(nil, nil).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	appInstance, err := newApp(s, "http://localhost:0007", "")
	if err != nil {
		assert.Error(t, err)
	}

	handler := appInstance.GetAuth()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		want   want
		method string
	}{
		{
			name: "content is empty",
			want: want{
				code:        http.StatusNoContent,
				contentType: "",
			},

			method: http.MethodGet,
		},
		{
			name: "incorrect method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, "/api/user/urls", nil)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func Test_deleteURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := storeMock.NewMockStore(ctrl)

	//установим условие: при любом вызове метода DeleteURLs не возвращались ошибки
	s.EXPECT().DeleteURLs(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	appInstance, err := newApp(s, "http://localhost:0007", "")
	if err != nil {
		assert.Error(t, err)
	}

	handler := appInstance.DeleteURLs()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	body := `["1", "2", "3"]`

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
				code:        http.StatusAccepted,
				contentType: "application/json",
			},
			body:   body,
			method: http.MethodDelete,
		},
		{
			name: "body is empty",
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "",
			},
			method: http.MethodDelete,
		},
		{
			name: "incorrect method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			body:   body,
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, "/api/user/urls", strings.NewReader(test.body))
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

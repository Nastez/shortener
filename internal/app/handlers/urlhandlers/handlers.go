package urlhandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Nastez/shortener/internal/app/models"
	"github.com/Nastez/shortener/internal/logger"
	"github.com/Nastez/shortener/internal/storage"
	"github.com/Nastez/shortener/utils"
)

type URLHandler struct {
	storage  storage.URLStorage
	baseAddr string
}

func New(storage storage.URLStorage, baseAddr string) (*URLHandler, error) {
	if storage == nil {
		return nil, errors.New("storage is empty")
	}

	if baseAddr == "" {
		return nil, errors.New("baseAddr is empty")
	}

	return &URLHandler{
		storage:  storage,
		baseAddr: baseAddr,
	}, nil
}

func (h *URLHandler) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var shortURL string

		if req.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info("can't read body")
			return
		}

		originalURL := string(body)
		if originalURL == "" {
			http.Error(w, "URL is empty", http.StatusBadRequest)
			return
		}

		defer req.Body.Close()

		generatedID := utils.GenerateID()
		h.storage.Save(originalURL, generatedID)

		shortURL = h.baseAddr + "/" + generatedID

		// устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "text/plain")
		// устанавливаем код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte(shortURL))
	}
}

func (h *URLHandler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		urlID := chi.URLParam(req, "id")
		if urlID == "" {
			http.Error(w, "urlID is missed", http.StatusBadRequest)
			return
		}

		if req.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var originalURL = h.storage.Get(urlID)
		fmt.Println("originalURL", originalURL)

		// устанавливаем заголовок Location
		w.Header().Set("Location", originalURL)
		// устанавливаем код 307
		w.WriteHeader(http.StatusOK)
	}
}

func (h *URLHandler) ShortenerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("got request with bad method", zap.String("method", req.Method))
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// десериализуем запрос в структуру модели
		logger.Log.Info("decoding request")
		var request models.Request
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&request); err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var shortURL string
		originalURL := request.URL

		generatedID := utils.GenerateID()
		h.storage.Save(originalURL, generatedID)
		shortURL = h.baseAddr + "/" + generatedID

		// заполняем модель ответа
		resp := models.Response{
			Result: shortURL,
		}

		// устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// устанавливаем код 201
		w.WriteHeader(http.StatusCreated)

		// сериализуем ответ сервера
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			logger.Log.Info("error encoding response", zap.Error(err))
			return
		}
		logger.Log.Info("sending HTTP 201 response")
	}
}

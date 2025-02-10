package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Nastez/shortener/internal/app/models"
	"github.com/Nastez/shortener/internal/logger"
	"github.com/Nastez/shortener/internal/store"
	"github.com/Nastez/shortener/utils"
)

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
	store                     store.Store
	baseAddr                  string
	databaseConnectionAddress string
}

// newApp принимает на вход внешние зависимости приложения и возвращает новый объект app
func newApp(s store.Store, baseAddr string, databaseConnectionAddress string) (*app, error) {
	if s == nil {
		return nil, errors.New("storage is empty")
	}

	if baseAddr == "" {
		return nil, errors.New("baseAddr is empty")
	}

	return &app{store: s, baseAddr: baseAddr, databaseConnectionAddress: databaseConnectionAddress}, nil
}

// GetPing проверяет соединение с базой данных
func (a *app) GetPing() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("got request with bad method", zap.String("method", req.Method))
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		if a == nil || a.databaseConnectionAddress == "" {
			logger.Log.Info("databaseConnectionAddress is nil")
			return
		}

		db, err := sql.Open("pgx", a.databaseConnectionAddress)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err = db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (a *app) ShortenerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

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
		shortURL = a.baseAddr + "/" + generatedID

		oldShortURL, err := a.store.Save(ctx, store.URL{
			OriginalURL: originalURL,
			ShortURL:    shortURL,
			GeneratedID: generatedID,
		})
		// наличие неспецифичной ошибки
		if err != nil && !errors.Is(err, store.ErrConflict) {
			logger.Log.Debug("cannot save urls in the store", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var resp models.Response

		if errors.Is(err, store.ErrConflict) {
			// ошибка специфична
			if oldShortURL == "" {
				logger.Log.Warn("oldShortURL is empty")
			}

			// заполняем модель ответа
			resp = models.Response{
				Result: oldShortURL,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			// устанавливаем код 409
			w.WriteHeader(http.StatusConflict)
			// сериализуем ответ сервера
			enc := json.NewEncoder(w)
			if err = enc.Encode(resp); err != nil {
				logger.Log.Info("error encoding response", zap.Error(err))
				return
			}
			logger.Log.Info("sending HTTP 409 response")
			return
		} else if err == nil {
			// заполняем модель ответа
			resp = models.Response{
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
}

func (a *app) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		urlID := chi.URLParam(req, "id")
		if urlID == "" {
			http.Error(w, "urlID is missed", http.StatusBadRequest)
			return
		}

		if req.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		//var originalURL = a.store.Get(urlID)
		originalURL, err := a.store.Get(ctx, urlID)
		if err != nil {
			fmt.Println(err)
			logger.Log.Debug("cannot get originalURL", zap.String("originalURL", originalURL), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("originalURL", originalURL)

		// устанавливаем заголовок Location
		w.Header().Set("Location", originalURL)
		fmt.Println("test originalURL", originalURL)
		// устанавливаем код 307
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (a *app) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var shortURL string
		ctx := req.Context()

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
		shortURL = a.baseAddr + "/" + generatedID

		if a == nil {
			return
		}

		oldShortURL, err := a.store.Save(ctx, store.URL{
			OriginalURL: originalURL,
			ShortURL:    shortURL,
			GeneratedID: generatedID,
		})

		// наличие неспецифичной ошибки
		if err != nil && !errors.Is(err, store.ErrConflict) {
			logger.Log.Debug("cannot save urls in the store", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if errors.Is(err, store.ErrConflict) {
			// ошибка специфична
			// устанавливаем заголовок Content-Type
			w.Header().Set("Content-Type", "text/plain")
			// устанавливаем код 409
			w.WriteHeader(http.StatusConflict)
			// пишем старый короткий url в тело ответа
			w.Write([]byte(oldShortURL))
			return
		}
		// устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "text/plain")
		// устанавливаем код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte(shortURL))

	}
}

func (a *app) PostBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		if req.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// десериализуем запрос в структуру модели
		logger.Log.Info("decoding request")
		var requestBatch models.PayloadBatch
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&requestBatch); err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var responseBatch models.ResponseBodyBatch

		for _, request := range requestBatch {
			var response = models.ResponseBatch{
				CorrelationID: request.CorrelationID,
				ShortURL:      a.baseAddr + "/" + request.CorrelationID,
			}
			responseBatch = append(responseBatch, response)
		}

		if len(responseBatch) > 0 {
			err := a.store.SaveBatch(ctx, requestBatch, responseBatch)
			if err != nil {
				logger.Log.Info("can't save batch in store")
				fmt.Println(err)
				return
			}
		} else {
			logger.Log.Info("responseBatch is empty")
		}

		// устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// устанавливаем код 201
		w.WriteHeader(http.StatusCreated)

		// сериализуем ответ сервера
		enc := json.NewEncoder(w)
		if err := enc.Encode(responseBatch); err != nil {
			logger.Log.Info("error encoding response", zap.Error(err))
			return
		}
		logger.Log.Info("sending HTTP 201 response")
	}
}

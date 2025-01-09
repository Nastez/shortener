package urlhandlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Nastez/shortener/internal/storage"
	"github.com/Nastez/shortener/utils"
)

type URLHandler struct {
	storage  storage.URLStorage
	baseAddr string
}

func New(storage storage.URLStorage, baseAddr string) (*URLHandler, error) {
	if storage != nil {
		return nil, errors.New("storage is empty")
	}

	if baseAddr != "" {
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
			log.Fatalln(err)
		}

		originalURL := string(body)
		if originalURL == "" {
			http.Error(w, "URL is empty", http.StatusBadRequest)
			return
		}

		defer req.Body.Close()

		generatedID := utils.GenerateID()
		h.storage.SaveOriginalURL(originalURL, generatedID)

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

		var originalURL = h.storage.GetOriginalURL(urlID)
		fmt.Println("originalURL", originalURL)

		// устанавливаем заголовок Location
		w.Header().Set("Location", originalURL)
		// устанавливаем код 307
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

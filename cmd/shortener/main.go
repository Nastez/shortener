package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var storeURL = make(map[string]string)

type ShortenerHandler struct{}

func main() {
	r := chi.NewRouter()

	r.Mount("/", ShortenerRoutes())

	http.ListenAndServe(":8080", r)
}

func ShortenerRoutes() chi.Router {
	r := chi.NewRouter()
	shortenerHandler := ShortenerHandler{}

	r.Post("/", shortenerHandler.postHandler(storeURL))
	r.Get("/{id}", shortenerHandler.getHandler)

	return r
}

func (s ShortenerHandler) postHandler(storeURL map[string]string) http.HandlerFunc {
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

		generatedID := generateID()
		storeURL[generatedID] = originalURL

		shortURL = "http://localhost:8080/" + generatedID

		// устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "text/plain")
		// устанавливаем код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte(shortURL))
	}
}

func (s ShortenerHandler) getHandler(w http.ResponseWriter, req *http.Request) {
	urlID := chi.URLParam(req, "id")
	if urlID == "" {
		http.Error(w, "urlID is missed", http.StatusBadRequest)
		return
	}

	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var originalURL = storeURL[urlID]

	// устанавливаем заголовок Location
	w.Header().Set("Location", originalURL)
	// устанавливаем код 307
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

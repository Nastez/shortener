package main

import (
	"fmt"
	"github.com/Nastez/shortener/config"
	"github.com/Nastez/shortener/utils"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//Задание по треку «Сервис сокращения URL»

//Добавьте возможность конфигурировать сервис с помощью переменных окружения:
//Адрес запуска HTTP-сервера — с помощью переменной SERVER_ADDRESS.
//Базовый адрес результирующего сокращённого URL — с помощью переменной BASE_URL.

//Приоритет параметров сервера должен быть таким:
//Если указана переменная окружения, то используется она.
//Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
//Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.

var storeURL = make(map[string]string)

type ShortenerHandler struct{}

func main() {
	config.ParseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()

	r.Mount("/", ShortenerRoutes(config.FlagBaseAddr))

	return http.ListenAndServe(":"+config.Port, r)
}

func ShortenerRoutes(baseAddr string) chi.Router {
	r := chi.NewRouter()
	shortenerHandler := ShortenerHandler{}

	r.Post("/", shortenerHandler.postHandler(storeURL, baseAddr))
	r.Get("/{id}", shortenerHandler.getHandler)

	return r
}

func (s ShortenerHandler) postHandler(storeURL map[string]string, baseAddr string) http.HandlerFunc {
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
		storeURL[generatedID] = originalURL

		if baseAddr == "http://localhost:" {
			http.Error(w, "port is empty", http.StatusBadRequest)
			return
		}

		shortURL = baseAddr + "/" + generatedID

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
	fmt.Println("originalURL", originalURL)

	// устанавливаем заголовок Location
	w.Header().Set("Location", originalURL)
	// устанавливаем код 307
	w.WriteHeader(http.StatusTemporaryRedirect)
}

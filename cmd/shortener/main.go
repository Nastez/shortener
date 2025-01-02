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

var storeURL = make(map[string]string)

type ShortenerHandler struct{}

func main() {
	config.ParseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", config.FlagRunAddr)
	fmt.Println("Running server on", config.FlagBaseAddr)

	r := chi.NewRouter()

	r.Mount("/", ShortenerRoutes(config.FlagBaseAddr))

	fmt.Println("portTestik", config.PortTest)
	return http.ListenAndServe(":"+config.PortTest, r)
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

		fmt.Println("req body is:", string(body))
		fmt.Println("req.URL.Path is:", req.URL.Path)
		fmt.Println("req.URL.RawQuery is:", req.URL.RawQuery)
		fmt.Println("req.Host is:", req.Host)
		fmt.Println("req.URL.Scheme is:", req.URL.Scheme)
		fmt.Println("req.URL.RequestURI() is:", req.URL.RequestURI())

		defer req.Body.Close()

		generatedID := utils.GenerateID()
		storeURL[generatedID] = originalURL

		if baseAddr == "http://localhost:" {
			http.Error(w, "port is empty", http.StatusBadRequest)
			return
		}
		fmt.Println("base addr:", baseAddr)

		shortURL = baseAddr + "/" + generatedID
		fmt.Println("shortURL is:", shortURL)

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

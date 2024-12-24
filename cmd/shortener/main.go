package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
)

var storeURL = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postHandler(storeURL))
	mux.HandleFunc(`/{id}`, getHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func postHandler(storeURL map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var shortURL string

		if req.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusBadRequest)
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

func getHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusBadRequest)
		return
	}

	var path = req.URL.Path
	var id = path[1:]
	var originalURL = storeURL[id]

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

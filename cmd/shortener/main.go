package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var originalURL = ""
var shortURL = "http://localhost:8080/EwHXdJfB "

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postHandler)
	mux.HandleFunc(`/{id}`, getHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func postHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Request Body:", string(body))
	defer req.Body.Close()

	// устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "text/plain")
	// устанавливаем код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte(shortURL))
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusBadRequest)
		return
	}

	var body string
	if err := req.ParseForm(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	for k, v := range req.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	// устанавливаем заголовок Location
	w.Header().Set("Location", originalURL)
	// устанавливаем код 307
	w.WriteHeader(http.StatusTemporaryRedirect)
	// пишем тело ответа
	w.Write([]byte(body))
}

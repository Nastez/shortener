package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

//Сервер должен быть доступен по адресу http://localhost:8080 и предоставлять два эндпоинта:
//Эндпоинт с методом POST и путём /. Сервер принимает в теле запроса строку URL как text/plain
//и возвращает ответ с кодом 201 и сокращённым URL как text/plain.

//Эндпоинт с методом GET и путём /{id}, где id — идентификатор сокращённого URL (например, /EwHXdJfB).
//	В случае успешной обработки запроса сервер возвращает ответ с кодом 307 и оригинальным URL в
//HTTP-заголовке Location.

//На любой некорректный запрос сервер должен возвращать ответ с кодом 400.

var originalUrl = ""
var shortUrl = "http://localhost:8080/EwHXdJfB "

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
	w.Write([]byte(shortUrl))
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
	w.Header().Set("Location", originalUrl)
	// устанавливаем код 307
	w.WriteHeader(http.StatusTemporaryRedirect)
	// пишем тело ответа
	w.Write([]byte(body))
}

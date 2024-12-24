package post_handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var shortURL = "http://localhost:8080/EwHXdJfB "

func PostHandler(w http.ResponseWriter, req *http.Request) {
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

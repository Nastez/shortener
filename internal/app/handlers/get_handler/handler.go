package get_handler

import (
	"fmt"
	"net/http"
)

var originalURL = ""

func GetHandler(w http.ResponseWriter, req *http.Request) {
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

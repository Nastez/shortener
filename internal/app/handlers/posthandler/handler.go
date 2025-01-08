package posthandler

import (
	"io"
	"log"
	"net/http"

	"github.com/Nastez/shortener/internal/storage"
	"github.com/Nastez/shortener/utils"
)

func PostHandler(s storage.URLStorage, baseAddr string) http.HandlerFunc {
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
		s.SaveOriginalURL(originalURL, generatedID)

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

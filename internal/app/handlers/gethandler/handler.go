package gethandler

import (
	"fmt"
	"net/http"

	"github.com/Nastez/shortener/internal/storage"

	"github.com/go-chi/chi/v5"
)

func GetHandler(s storage.URLStorage) http.HandlerFunc {
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

		//var originalURL = storage.StoreURL[urlID]
		var originalURL = s.GetOriginalURL(urlID)
		fmt.Println("originalURL", originalURL)

		// устанавливаем заголовок Location
		w.Header().Set("Location", originalURL)
		// устанавливаем код 307
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

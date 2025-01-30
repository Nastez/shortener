package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/Nastez/shortener/config"
	"github.com/Nastez/shortener/internal/app/handlers/urlhandlers"
	"github.com/Nastez/shortener/internal/logger"
	"github.com/Nastez/shortener/internal/saver"
	"github.com/Nastez/shortener/internal/storage"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	if err = run(cfg); err != nil {
		panic(err)
	}
}

func run(cfg *config.Config) error {
	r := chi.NewRouter()

	defer os.Remove(cfg.FileName)

	err := saver.SaveFile(cfg.FileName)
	if err != nil {
		return err
	}

	routes, err := ShortenerRoutes(cfg.BaseURL)
	if err != nil {
		return err
	}

	r.Mount("/", routes)

	return http.ListenAndServe(":"+cfg.Port, r)
}

func ShortenerRoutes(baseAddr string) (chi.Router, error) {
	r := chi.NewRouter()

	storeURL := storage.MemoryStorage{}

	if baseAddr == "http://localhost:" {
		return nil, errors.New("port is empty")
	}

	handlers, err := urlhandlers.New(storeURL, baseAddr)
	if err != nil {
		return nil, err
	}

	r.Post("/", logger.WithLogging(GzipMiddleware(handlers.PostHandler())))
	r.Get("/{id}", logger.WithLogging(GzipMiddleware(handlers.GetHandler())))
	r.Post("/api/shorten", logger.WithLogging(GzipMiddleware(handlers.ShortenerHandler())))

	return r, nil
}

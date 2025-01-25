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
	err := config.ParseFlagsAndEnv()
	if err != nil {
		log.Fatalln(err)
	}

	if err = run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()

	fileName := &config.FileName
	defer os.Remove(*fileName)

	err := saver.SaveFile(*fileName)
	if err != nil {
		return err
	}

	routes, err := ShortenerRoutes(config.FlagBaseAddr)
	if err != nil {
		return err
	}

	r.Mount("/", routes)

	return http.ListenAndServe(":"+config.Port, r)
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

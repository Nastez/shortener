package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Nastez/shortener/config"
	getHandler "github.com/Nastez/shortener/internal/app/handlers/gethandler"
	postHandler "github.com/Nastez/shortener/internal/app/handlers/posthandler"
	"github.com/Nastez/shortener/internal/storage"
)

func main() {
	err := config.ParseFlags()
	if err != nil {
		log.Fatalln(err)
	}

	if err = run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()

	r.Mount("/", ShortenerRoutes(config.FlagBaseAddr))

	return http.ListenAndServe(":"+config.Port, r)
}

func ShortenerRoutes(baseAddr string) chi.Router {
	r := chi.NewRouter()

	storeURL := storage.MemoryStorage{}

	r.Post("/", postHandler.PostHandler(storeURL, baseAddr))
	r.Get("/{id}", getHandler.GetHandler(storeURL))

	return r
}

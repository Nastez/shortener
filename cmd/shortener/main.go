package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/Nastez/shortener/config"
	"github.com/Nastez/shortener/internal/logger"
	"github.com/Nastez/shortener/internal/saver"
	"github.com/Nastez/shortener/internal/storage"
	"github.com/Nastez/shortener/internal/store/pg"
	"github.com/Nastez/shortener/internal/storeconfig"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	if cfg.DatabaseConnectionAddress != "" {
		// создаём соединение с СУБД PostgreSQL с помощью аргумента командной строки
		conn, err := sql.Open("pgx", cfg.DatabaseConnectionAddress)
		if err != nil {
			return err
		}

		// создаём экземпляр приложения, передавая реализацию хранилища pg в качестве внешней зависимости
		appInstance, err := newApp(pg.NewStore(conn), cfg.BaseURL, cfg.DatabaseConnectionAddress)
		if err != nil {
			return err
		}
		storeconfig.NewStoreConfig(conn).Bootstrap(context.Background())

		routes, err := ShortenerRoutes(cfg.BaseURL, *appInstance)
		if err != nil {
			return err
		}

		r.Mount("/", routes)
	} else if cfg.DatabaseConnectionAddress == "" && cfg.FileName != "events.log" {
		defer os.Remove(cfg.FileName)

		err := saver.SaveFile(cfg.FileName)
		if err != nil {
			return err
		}

		appInstance, err := newApp(storage.New(), cfg.BaseURL, cfg.DatabaseConnectionAddress)
		if err != nil {
			return err
		}

		routes, err := ShortenerRoutes(cfg.BaseURL, *appInstance)
		if err != nil {
			return err
		}

		r.Mount("/", routes)
	} else if cfg.DatabaseConnectionAddress == "" && cfg.FileName == "events.log" {
		appInstance, err := newApp(storage.New(), cfg.BaseURL, cfg.DatabaseConnectionAddress)
		if err != nil {
			return err
		}

		routes, err := ShortenerRoutes(cfg.BaseURL, *appInstance)
		if err != nil {
			return err
		}

		r.Mount("/", routes)

	} else {
		logger.Log.Error("can't run app")
	}

	return http.ListenAndServe(":"+cfg.Port, r)
}

func ShortenerRoutes(baseAddr string, appInstance app) (chi.Router, error) {
	r := chi.NewRouter()

	if baseAddr == "http://localhost:" {
		return nil, errors.New("port is empty")
	}

	r.Post("/", logger.WithLogging(GzipMiddleware(appInstance.PostHandler())))
	r.Get("/{id}", logger.WithLogging(GzipMiddleware(appInstance.GetHandler())))
	r.Post("/api/shorten", logger.WithLogging(GzipMiddleware(appInstance.ShortenerHandler())))
	r.Get("/ping", logger.WithLogging(GzipMiddleware(appInstance.GetPing())))
	r.Post("/api/shorten/batch", logger.WithLogging(GzipMiddleware(appInstance.PostBatch())))
	r.Get("/api/user/urls", logger.WithLogging(GzipMiddleware(appInstance.GetAuth())))
	r.Delete("/api/user/urls", logger.WithLogging(GzipMiddleware(appInstance.DeleteURLs())))

	return r, nil
}

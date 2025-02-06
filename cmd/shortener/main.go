package main

import (
	"database/sql"
	"errors"
	"fmt"
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

	// создаём соединение с СУБД PostgreSQL с помощью аргумента командной строки
	conn, err := sql.Open("pgx", cfg.DatabaseConnectionAddress)
	if err != nil {
		return err
	}

	// Проверка соединения
	err = conn.Ping()
	if err != nil {
		log.Fatal("ошибка подключения:", err)
	}
	fmt.Println("подключение к БД успешно")

	// создаём экземпляр приложения, передавая реализацию хранилища pg в качестве внешней зависимости
	// appInstance := newApp(pg.NewStore(conn), cfg.BaseURL, cfg.DatabaseURI)
	//appInstance, err := urlhandlers.New(pg.NewStore(conn), cfg.BaseURL, cfg.DatabaseConnectionAddress)
	//if err != nil {
	//	return err
	//}

	routes, err := ShortenerRoutes(cfg.BaseURL, cfg.DatabaseConnectionAddress)
	if err != nil {
		return err
	}

	r.Mount("/", routes)

	return http.ListenAndServe(":"+cfg.Port, r)
}

func ShortenerRoutes(baseAddr string, databaseConnectionAddress string) (chi.Router, error) {
	r := chi.NewRouter()

	storeURL := storage.MemoryStorage{}

	if baseAddr == "http://localhost:" {
		return nil, errors.New("port is empty")
	}

	if databaseConnectionAddress == "" {
		return nil, errors.New("get databaseConnectionAddress error")
	}

	handlers, err := urlhandlers.New(storeURL, baseAddr, databaseConnectionAddress)
	if err != nil {
		return nil, err
	}

	r.Post("/", logger.WithLogging(GzipMiddleware(handlers.PostHandler())))
	r.Get("/{id}", logger.WithLogging(GzipMiddleware(handlers.GetHandler())))
	r.Post("/api/shorten", logger.WithLogging(GzipMiddleware(handlers.ShortenerHandler())))
	r.Get("/ping", logger.WithLogging(GzipMiddleware(handlers.GetPing())))

	return r, nil
}

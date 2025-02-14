package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v6"
)

// Env Переменные окружения
type Env struct {
	ServerAddress             string `env:"SERVER_ADDRESS"`
	BaseURL                   string `env:"BASE_URL"`
	FileStoragePath           string `env:"FILE_STORAGE_PATH"`
	DatabaseConnectionAddress string `env:"DATABASE_DSN"`
}

type Config struct {
	ServerAddress             string
	BaseURL                   string
	FileStoragePath           string
	Port                      string
	FileName                  string
	DatabaseConnectionAddress string
}

// New обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func New() (*Config, error) {
	var envConf Env
	err := env.Parse(&envConf)
	if err != nil {
		return nil, errors.New("can't parse env")
	}

	log.Println(envConf)

	var (
		serverAddress             string
		baseURL                   string
		fileStoragePath           string
		port                      string
		fileName                  string
		databaseConnectionAddress string
	)

	psDefault := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `shortener`, `k8tego`, `shortener`)

	flag.StringVar(&serverAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base address before a short URL")
	flag.StringVar(&fileStoragePath, "f", "events.log", "file storage path")
	flag.StringVar(&databaseConnectionAddress, "d", psDefault, "database connection address")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envConf.ServerAddress != "" {
		serverAddress = envConf.ServerAddress
	}

	if envConf.BaseURL != "" {
		baseURL = envConf.BaseURL
	}

	if fileStoragePath != "" {
		fileName = fileStoragePath
	}

	if envConf.FileStoragePath != "" {
		fileName = envConf.FileStoragePath
	}

	if envConf.FileStoragePath == "" && fileStoragePath != "" {
		fileName = fileStoragePath
	}

	if envConf.DatabaseConnectionAddress != "" {
		databaseConnectionAddress = envConf.DatabaseConnectionAddress
	}

	if baseURL == "http://localhost:" || baseURL == "http://localhost:/" {
		fmt.Fprintf(os.Stderr, "Invalid base address: %s (must has format http://localhost:8080/)\n", baseURL)
		os.Exit(1)
	}

	port = splitRunAddr(serverAddress)

	if !validatePort(port) {
		return nil, errors.New("invalid port number")
	}

	return &Config{
		ServerAddress:             serverAddress,
		BaseURL:                   baseURL,
		FileStoragePath:           fileStoragePath,
		Port:                      port,
		FileName:                  fileName,
		DatabaseConnectionAddress: databaseConnectionAddress,
	}, nil
}

func validatePort(port string) bool {
	match, _ := regexp.MatchString(`^[0-9]+$`, port)
	return match
}

func splitRunAddr(runAddr string) string {
	splittedRunAddr := strings.Split(runAddr, ":")
	port := splittedRunAddr[1]

	return port
}

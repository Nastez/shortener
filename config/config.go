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

var FlagRunAddr string
var FlagBaseAddr string
var FlagFileStoragePath string
var Port string
var FileName string

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

// ParseFlagsAndEnv обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlagsAndEnv() error {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return errors.New("can't parse env")
	}

	log.Println(cfg)

	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8080", "base address before a short URL")
	flag.StringVar(&FlagFileStoragePath, "f", "events.log", "file storage path")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if cfg.ServerAddress != "" {
		FlagRunAddr = cfg.ServerAddress
	}

	if cfg.BaseURL != "" {
		FlagBaseAddr = cfg.BaseURL
	}

	if FlagFileStoragePath != "" {
		FileName = FlagFileStoragePath
	}

	if cfg.FileStoragePath != "" {
		FileName = cfg.FileStoragePath
	}

	if cfg.FileStoragePath == "" && FlagFileStoragePath != "" {
		FileName = FlagFileStoragePath
	}

	if FlagBaseAddr == "http://localhost:" || FlagBaseAddr == "http://localhost:/" {
		fmt.Fprintf(os.Stderr, "Invalid base address: %s (must has format http://localhost:8080/)\n", FlagBaseAddr)
		os.Exit(1)
	}

	port := splitRunAddr(FlagRunAddr)

	if !validatePort(port) {
		return errors.New("invalid port number")
	}

	Port = port

	return err
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

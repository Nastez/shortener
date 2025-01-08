package config

import (
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
var Port string

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() error {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8080", "base address before a short URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if cfg.ServerAddress != "" {
		FlagRunAddr = cfg.ServerAddress
	}

	if cfg.BaseURL != "" {
		FlagBaseAddr = cfg.BaseURL
	}

	if FlagBaseAddr == "http://localhost:" || FlagBaseAddr == "http://localhost:/" {
		fmt.Fprintf(os.Stderr, "Invalid base address: %s (must has format http://localhost:8080/)\n", FlagBaseAddr)
		os.Exit(1)
	}

	port := splitRunAddr(FlagRunAddr)

	if !validatePort(port) {
		log.Fatalf("Invalid port number: %s", port)
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

package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var invalidPort string

var FlagRunAddr int
var FlagBaseAddr string

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {

	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.IntVar(&FlagRunAddr, "a", 8080, "address and port to run server")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8080/", "base address before a short URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if FlagBaseAddr == "http://localhost:" || FlagBaseAddr == "http://localhost:/" {
		fmt.Fprintf(os.Stderr, "Invalid base address: %s (must have format http://localhost:8080/)\n", FlagBaseAddr)
		os.Exit(1)
	}
	if !validatePort(FlagRunAddr) {
		log.Fatalf("Invalid port number: %d", FlagRunAddr)
	}
}

func validatePort(port int) bool {
	return port > 0 && port <= 65535
}

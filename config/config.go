package config

import (
	"flag"
	"fmt"
	"os"
)

var invalidPort string

var FlagRunAddr string
var FlagBaseAddr string

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {

	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8080/", "base address before a short URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if FlagBaseAddr == "http://localhost:" || FlagBaseAddr == "http://localhost:/" {
		fmt.Fprintf(os.Stderr, "Invalid base address: %s (must have format http://localhost:8080/)\n", FlagBaseAddr)
		os.Exit(1)
	}

	if len([]rune(FlagRunAddr)) > 5 {
		fmt.Fprintf(os.Stderr, "Invalid port: %s (port must have format :8080)\n", FlagBaseAddr)
		os.Exit(1)
	}
}

package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/Nastez/shortener/utils"
)

var invalidPort = utils.GenerateID()

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

	if FlagBaseAddr == "http://localhost:" || FlagBaseAddr == "http://localhost:/" || FlagBaseAddr == "http://localhost:"+invalidPort || FlagBaseAddr == "http://localhost:/"+invalidPort {
		fmt.Fprintf(os.Stderr, "Invalid base address: %s (must have format http://localhost:8080/)\n", FlagBaseAddr)
		os.Exit(1)
	}
}

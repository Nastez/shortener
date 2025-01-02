package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var FlagRunAddr string
var FlagBaseAddr string
var PortTest string

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8080", "base address before a short URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if FlagBaseAddr == "http://localhost:" || FlagBaseAddr == "http://localhost:/" {
		fmt.Fprintf(os.Stderr, "Invalid base address: %s (must has format http://localhost:8080/)\n", FlagBaseAddr)
		os.Exit(1)
	}

	port := splitRunAddr(FlagRunAddr)

	if !validatePort(port) {
		log.Fatalf("Invalid port number: %s", port)
	}

	PortTest = port
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

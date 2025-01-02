package config

import "flag"

type Flags struct {
	FlagRunAddr  string
	FlagBaseAddr string
}

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&InitFlags().FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&InitFlags().FlagBaseAddr, "b", "http://localhost:8080/", "base address before a short URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}

func InitFlags() *Flags {
	var f *Flags
	f = new(Flags)
	f.FlagRunAddr = ":8080"
	f.FlagBaseAddr = "http://localhost:8080/"

	return f
}

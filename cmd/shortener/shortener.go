package main

import (
	"flag"
	"link-shortener/internal/shortener"
	"log"
)

func main() {
	var storageType string
	flag.StringVar(&storageType, "storage", "memory", "Тип хранилища")
	flag.Parse()

	if err := shortener.Run(storageType); err != nil {
		log.Fatalf("Ошибка запуска сервиса: %v", err)
	}
}

package storage

import "time"

const storageTime = 24 * 7 * time.Hour

type Storage interface {
	// ShortenAndSaveUrl сохраняет оригинальный shortUrl и возвращает укороченный
	ShortenAndSaveUrl(url string) (string, error)

	// GetUrl возвращает оригинальный shortUrl по укороченному
	GetUrl(shortUrl string) (string, error)

	// Close завершает работу хранилища
	Close() error
}

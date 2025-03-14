package storage

import "time"

const storageTime = 24 * 7 * time.Hour

type Storage interface {
	// ShortenAndSaveUrl сохраняет оригинальный Url и возвращает укороченный
	ShortenAndSaveUrl(url string) (string, error)

	// GetUrl возвращает оригинальный Url по укороченному
	GetUrl(shortUrl string) (string, error)

	// Close завершает работу хранилища
	Close() error
}

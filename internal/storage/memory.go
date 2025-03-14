package storage

import (
	"fmt"
	slg "link-shortener/internal/short_link_generator"
	"sync"
	"time"
)

// shortExpirationPair хранит сокращенную ссылку и дату когда заканчивается ее "срок действия"
type shortExpirationPair struct {
	shortUrl   string
	expireDate time.Time
}

// MemoryStorage - in-memory хранилище
type MemoryStorage struct {
	sync.RWMutex
	shortToOrig  map[string]string
	origToShort  map[string]shortExpirationPair
	cleaningDone chan struct{}
}

// NewMemoryStorage создает новое in-memory хранилищу
func NewMemoryStorage() Storage {
	storage := &MemoryStorage{
		shortToOrig:  make(map[string]string),
		origToShort:  make(map[string]shortExpirationPair),
		cleaningDone: make(chan struct{}, 1),
	}
	go storage.cleaningLoop()
	return storage

}

// Close прекращает цикл удаления давно зарегестрированных ссылок
func (m *MemoryStorage) Close() error {
	m.cleaningDone <- struct{}{}
	return nil
}

// ShortenAndSaveUrl сокращает переданный Url и сохраняет сокращенный
func (m *MemoryStorage) ShortenAndSaveUrl(url string) (string, error) {
	m.Lock()
	defer m.Unlock()

	shortUrl := ""

	if _, ok := m.origToShort[url]; ok {
		shortUrl = m.origToShort[url].shortUrl
		m.origToShort[url] = shortExpirationPair{
			shortUrl:   shortUrl,
			expireDate: time.Now().Add(storageTime),
		}
		return shortUrl, nil
	}

	for attempts := 0; attempts < 25; attempts++ {
		shortUrl = slg.GenerateShortLink()

		if _, exists := m.shortToOrig[shortUrl]; !exists {
			m.shortToOrig[shortUrl] = url
			m.origToShort[url] = shortExpirationPair{
				shortUrl:   shortUrl,
				expireDate: time.Now().Add(storageTime),
			}
			return shortUrl, nil
		}

	}

	return "", fmt.Errorf("не удалось сгенерировать короткую ссылку")
}

// GetUrl возвращает оригинальную ссылку по короткой, если она есть
func (m *MemoryStorage) GetUrl(shortUrl string) (string, error) {
	m.Lock()
	defer m.Unlock()

	if origUrl, ok := m.shortToOrig[shortUrl]; ok {
		m.origToShort[origUrl] = shortExpirationPair{ // обновляем время хранения ссылки
			shortUrl:   shortUrl,
			expireDate: time.Now().Add(storageTime),
		}
		return m.shortToOrig[shortUrl], nil
	}
	return "", fmt.Errorf("не удалось найти короткую ссылку")
}

// cleaningLoop удаляет давно зарегестрированные ссылки
func (m *MemoryStorage) cleaningLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-m.cleaningDone:
			return
		case <-ticker.C:
			now := time.Now()
			m.Lock()
			for url, pair := range m.origToShort {
				if now.After(pair.expireDate) {
					delete(m.origToShort, url)
					delete(m.shortToOrig, pair.shortUrl)
				}
			}
			m.Unlock()
		}

	}
}

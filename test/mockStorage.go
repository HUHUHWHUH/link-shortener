package test

import "fmt"

const (
	shortUrlToSave       = "1234567890"
	shortUrlToInitialize = "abcdefghij"
	originalUrl          = "https://www.wikipedia.org/"
	nonExistentShortUrl  = "doesntExst"
	notValidUrl          = "httsp:/sdsd"
)

// mockStorage - мок-хранилище, реализующее интерфейс Storage
type mockStorage struct {
	data map[string]string
}

func (m *mockStorage) Close() error {
	return nil
}

func (m *mockStorage) ShortenAndSaveUrl(url string) (string, error) {
	m.data[url] = shortUrlToSave
	return shortUrlToSave, nil
}

func (m *mockStorage) GetUrl(shortUrl string) (string, error) {
	if orig, ok := m.data[shortUrl]; ok {
		return orig, nil
	}
	return "", fmt.Errorf("не удалось найти короткую ссылку")
}

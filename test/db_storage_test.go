package test

import (
	"errors"
	"link-shortener/internal/storage"
	"testing"
)

func TestNewDbStorage_InvalidConn(t *testing.T) {
	_, err := storage.NewDbStorage("invalid_connection_string")
	if err == nil {
		t.Error("ожидалась ошибка при создании DbStorage с некорректной строкой подключения")
	}
}

func TestDbStorage_ShortenAndSaveUrl(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	storeDb, err := storage.NewDbStorage(connStr)
	if err != nil {
		t.Errorf("ошибка создания бд: %v", err)
	}

	url := "https://www.wikipedia.org/"
	shortUrl, err := storeDb.ShortenAndSaveUrl(url)
	if err != nil {
		t.Fatalf("ошибка при сохранении ссылки: %v", err)
	}

	shortUrl2, err := storeDb.ShortenAndSaveUrl(url)
	if err != nil {
		t.Fatalf("ошибка при повторном сохранении ссылки: %v", err)
	}
	if shortUrl != shortUrl2 {
		t.Errorf("ожидалась одинаковая короткая ссылка, получено %s и %s", shortUrl, shortUrl2)
	}
}

func TestDbStorage_ShortenAndSaveAndGetUrl(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	storeDb, err := storage.NewDbStorage(connStr)
	if err != nil {
		t.Errorf("ошибка создания бд: %v", err)
	}
	url := "https://www.wikipedia.org/"

	shortUrl, err := storeDb.ShortenAndSaveUrl(url)
	if err != nil {
		t.Errorf("ошибка при сохранении ссылки: %v", err)
	}

	gotUrl, err := storeDb.GetUrl(shortUrl)
	if err != nil {
		t.Errorf("ошибка при получении ссылкиЖ %v", err)
	}

	if gotUrl != url {
		t.Errorf("ожидалось %s, получено %s", url, gotUrl)
	}
}

func TestDbStorage_GetNonExistentUrl(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	storeDb, err := storage.NewDbStorage(connStr)
	if err != nil {
		t.Errorf("ошибка создания бд: %v", err)
	}
	nonExistentUrl := "https://doesnt_exist"
	var errNotFound = errors.New("не удалось найти короткую ссылку: sql: no rows in result set")

	gotUrl, err := storeDb.GetUrl(nonExistentUrl)
	if gotUrl != "" {
		t.Errorf("ожидалась пустая строка, получено: %s", gotUrl)
	}
	if err != nil && err.Error() != errNotFound.Error() {
		t.Errorf("ожидалась ошибка: \"не удалось найти короткую ссылку\", получено: %v", err)
	}
}

package test

import (
	"errors"
	"link-shortener/internal/storage"
	"testing"
)

func TestMemoryStorage_ShortenAndSaveUrl(t *testing.T) {
	store := storage.NewMemoryStorage()

	url := "https://www.wikipedia.org/"
	shortUrl, err := store.ShortenAndSaveUrl(url)
	if err != nil {
		t.Fatalf("ошибка при сохранении ссылки: %v", err)
	}

	shortUrl2, err := store.ShortenAndSaveUrl(url)
	if err != nil {
		t.Fatalf("ошибка при повторном сохранении ссылк: %v", err)
	}
	if shortUrl != shortUrl2 {
		t.Errorf("ожидалась одинаковая короткая ссылка, получено %s и %s", shortUrl, shortUrl2)
	}

}

func TestMemoryStorage_SaveAndGetUrl(t *testing.T) {
	store := storage.NewMemoryStorage()
	url := "https://www.wikipedia.org/"

	shortUrl, err := store.ShortenAndSaveUrl(url)
	if err != nil {
		t.Errorf("ошибка при сохранении ссылки: %v", err)
	}

	gotUrl, err := store.GetUrl(shortUrl)
	if err != nil {
		t.Errorf("ошибка при получении ссылкиЖ %v", err)
	}

	if gotUrl != url {
		t.Errorf("ожидалось %s, получено %s", url, gotUrl)
	}
}

func TestMemoryStorage_GetNonExistentUrl(t *testing.T) {
	store := storage.NewMemoryStorage()
	nonExistentUrl := "https://doesnt_exist"
	var errNotFound = errors.New("не удалось найти короткую ссылку")

	gotUrl, err := store.GetUrl(nonExistentUrl)
	if gotUrl != "" {
		t.Errorf("ожидалась пустая строка, получено: %s", gotUrl)
	}
	if err != nil && err.Error() != errNotFound.Error() {
		t.Errorf("ожидалась ошибка: \"не удалось найти короткую ссылку\", получено: %v", err)
	}

}

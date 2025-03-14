package test

import (
	"encoding/json"
	"link-shortener/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetFromShortUrl_Success(t *testing.T) {
	mockSt := &mockStorage{data: map[string]string{shortUrlToInitialize: originalUrl}}

	h := handlers.NewHandler(mockSt)

	req := httptest.NewRequest(http.MethodGet, "/"+shortUrlToInitialize, nil)
	response := httptest.NewRecorder()

	h.HandleGetFromShortUrl(response, req)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("ожидался статус 200, получен %d", status)
	}

	var decodedResponse handlers.OriginalUrlResponse
	if err := json.NewDecoder(response.Body).Decode(&decodedResponse); err != nil {
		t.Errorf("ошибка декодирования тела ответа: %v", err)
	}

	if decodedResponse.Url != originalUrl {
		t.Errorf("ожидалась ссылка %s, получено %s", originalUrl, decodedResponse.Url)
	}
}

func TestHandleGetFromShortUrl_Fail(t *testing.T) {
	mockSt := &mockStorage{data: map[string]string{shortUrlToInitialize: originalUrl}}

	h := handlers.NewHandler(mockSt)

	req := httptest.NewRequest(http.MethodGet, "/"+nonExistentShortUrl, nil)
	response := httptest.NewRecorder()

	h.HandleGetFromShortUrl(response, req)

	if status := response.Code; status != http.StatusNotFound {
		t.Errorf("ожидался статус 404, получен %d", status)
	}
}

func TestHandleGetFromShortUrl_SendingIncorrectHttpMethod(t *testing.T) {
	mockSt := &mockStorage{data: map[string]string{shortUrlToInitialize: originalUrl}}

	h := handlers.NewHandler(mockSt)

	req := httptest.NewRequest(http.MethodPost, "/"+nonExistentShortUrl, nil)
	response := httptest.NewRecorder()

	h.HandleGetFromShortUrl(response, req)

	if status := response.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("ожидался статус 405, получен %d", status)
	}
}

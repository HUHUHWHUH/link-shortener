package test

import (
	"bytes"
	"encoding/json"
	"link-shortener/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleShortenUrl_Success(t *testing.T) {
	mockSt := &mockStorage{data: map[string]string{shortUrlToInitialize: originalUrl}}
	h := handlers.NewHandler(mockSt)

	postReq := handlers.PostRequest{Url: originalUrl}

	reqBody, err := json.Marshal(postReq)
	if err != nil {
		t.Fatalf("ошибка сериализации запроса: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(reqBody))
	response := httptest.NewRecorder()

	h.HandleShortenUrl(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("ожидался статус 200, получен %d", response.Code)
	}

	var resp handlers.ShorUrlResponse
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		t.Errorf("ошибка декодирования ответа: %v", err)
	}

	if resp.ShortUrl != shortUrlToSave {
		t.Errorf("ожидалась короткая ссылка %s, получено %s", shortUrlToSave, resp.ShortUrl)
	}
}

func TestHandleShortenUrl_SendingIncorrectHttpMethod(t *testing.T) {
	mockSt := &mockStorage{data: map[string]string{shortUrlToInitialize: originalUrl}}
	h := handlers.NewHandler(mockSt)

	postReq := handlers.PostRequest{Url: originalUrl}

	reqBody, err := json.Marshal(postReq)
	if err != nil {
		t.Fatalf("ошибка сериализации запроса: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/shorten", bytes.NewReader(reqBody))
	response := httptest.NewRecorder()

	h.HandleShortenUrl(response, req)

	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("ожидался статус 405, получен %d", response.Code)
	}
}

func TestHandleShortenUrl_SendingRequestWithIncorrectBody(t *testing.T) {
	mockSt := &mockStorage{data: map[string]string{shortUrlToInitialize: originalUrl}}
	h := handlers.NewHandler(mockSt)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader([]byte("incorrect body")))
	response := httptest.NewRecorder()

	h.HandleShortenUrl(response, req)

	if response.Code != http.StatusBadRequest {
		t.Errorf("ожидался статус 400, получен %d", response.Code)
	}
}

func TestHandleShortenUrl_SendingRequestWithNotValidUrl(t *testing.T) {
	mockSt := &mockStorage{data: map[string]string{shortUrlToInitialize: originalUrl}}
	h := handlers.NewHandler(mockSt)
	postReq := handlers.PostRequest{Url: notValidUrl}
	reqBody, err := json.Marshal(postReq)
	if err != nil {
		t.Fatalf("ошибка сериализации запроса: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(reqBody))
	response := httptest.NewRecorder()

	h.HandleShortenUrl(response, req)

	if response.Code != http.StatusBadRequest {
		t.Errorf("ожидался статус 400, получен %d", response.Code)
	}
}

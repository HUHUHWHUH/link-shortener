package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// isValidUrl проверяет, что строка является валидной ссылкой
func isValidUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// HandleShortenUrl обрабатывает POST запрос на создание короткой ссылки
func (h *Handler) HandleShortenUrl(writer http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodPost {
		http.Error(writer, "некорректный запрос", http.StatusMethodNotAllowed)
		return
	}

	var req PostRequest
	decoder := json.NewDecoder(httpRequest.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(writer, "ошибка при десериализации тела запроса", http.StatusBadRequest)
		return
	}

	if !isValidUrl(req.Url) {
		http.Error(writer, "некорректная ссылка", http.StatusBadRequest)
		return
	}

	shortUrl, err := h.storage.ShortenAndSaveUrl(req.Url)
	if err != nil {
		http.Error(writer, "ошибка при сохранении ссылки", http.StatusInternalServerError)
		return
	}

	resp := ShorUrlResponse{ShortUrl: shortUrl}
	writer.Header().Set("content-type", "application/json")
	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		http.Error(writer, fmt.Sprintf("ошибка при формировании ответа: %v", err), http.StatusInternalServerError)
		return
	}
}

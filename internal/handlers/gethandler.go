package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HandleGetFromShortUrl обрабатывает GET-запрос для получения оригинальной ссылки по короткой
func (h *Handler) HandleGetFromShortUrl(writer http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet {
		http.Error(writer, "некорректный запрос", http.StatusMethodNotAllowed)
		return
	}

	shortUrl := httpRequest.URL.Path[1:] // убираем / в начале

	url, err := h.storage.GetUrl(shortUrl)
	if err != nil {
		http.Error(writer, "ссылка не найдена", http.StatusNotFound)
		return
	}

	resp := OriginalUrlResponse{Url: url}
	writer.Header().Set("content-type", "application/json")
	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		http.Error(writer, fmt.Sprintf("ошибка при формировании ответа: %v", err), http.StatusInternalServerError)
		return
	}
}

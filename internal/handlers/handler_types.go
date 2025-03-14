package handlers

import "link-shortener/internal/storage"

// Handler - обработчик http запросов, содержит хранилище, заданное при запуске приложения
type Handler struct {
	storage storage.Storage
}

// NewHandler создает новый экземпляр Handler
func NewHandler(st storage.Storage) *Handler {
	return &Handler{storage: st}
}

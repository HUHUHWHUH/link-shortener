package shortener

import (
	"context"
	"errors"
	"fmt"
	"link-shortener/internal/handlers"
	"link-shortener/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	port = ":8000"
)

// Run - основная точка запуска сервиса
func Run(storageType string) error {
	store, err := createStorage(storageType)
	if err != nil {
		return fmt.Errorf("ошибка при создании хранилища: %w", err)
	}
	defer store.Close()

	h := handlers.NewHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.HandleShortenUrl)
	mux.HandleFunc("/", h.HandleGetFromShortUrl)

	server := &http.Server{ // создание сервера, что бы была возможность выполнить graceful shutdown
		Addr:    port,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка во время завершения работы сервера: %w", err)
	}
	return nil
}

// createStorage создает хранилище в зависимости от переданного типа хранилища
func createStorage(storageType string) (storage.Storage, error) {
	switch storageType {
	case "postgres":
		pgConnStr := os.Getenv("POSTGRES_CONN")
		db, err := storage.NewDbStorage(pgConnStr)
		return db, err
	case "memory":
		return storage.NewMemoryStorage(), nil
	default:
		return nil, fmt.Errorf("неправильный тип хранилища")
	}
}

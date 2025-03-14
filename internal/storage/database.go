package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"link-shortener/internal/short_link_generator"
	"log"
	"time"
)

// DbStorage хранилище с базой данных
type DbStorage struct {
	db               *sql.DB
	cleaningDone     chan struct{}
	cleaningFinished chan struct{}

	preparedSelectShort    *sql.Stmt
	preparedUpdate         *sql.Stmt
	preparedInsert         *sql.Stmt
	preparedDelete         *sql.Stmt
	preparedSelectOriginal *sql.Stmt
}

// closeStmt закрывает перееданный подготовленный запрос
func closeStmt(stmt *sql.Stmt) {
	if stmt != nil {
		stmt.Close()
	}
}

// closeAllStmts закрывает все подготовленные запросы
func (d *DbStorage) closeAllStmts() {
	closeStmt(d.preparedSelectOriginal)
	closeStmt(d.preparedDelete)
	closeStmt(d.preparedUpdate)
	closeStmt(d.preparedSelectShort)
	closeStmt(d.preparedInsert)
}

// Close закрывает базу данных и прекращает цикл удаления давно зарегестрированных ссылок
func (d *DbStorage) Close() error {
	d.cleaningDone <- struct{}{}

	select {
	case <-d.cleaningFinished:
	case <-time.After(5 * time.Second):
		log.Println("cleaningLoop не завершился во время")
	}

	d.closeAllStmts()
	return d.db.Close()
}

// NewDbStorage создает новое хранилище с бд
func NewDbStorage(connStr string) (Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть бд: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS links (
            original_url TEXT PRIMARY KEY,
            short_url CHAR(10) NOT NULL UNIQUE,
            expiration_date TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '7 days')
        );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании бд: %w", err)
	}
	storage := &DbStorage{
		db:               db,
		cleaningDone:     make(chan struct{}, 1),
		cleaningFinished: make(chan struct{}, 1),
	}

	storage.preparedSelectShort, err = db.Prepare("SELECT short_url FROM links WHERE original_url = $1")
	if err != nil {
		storage.closeAllStmts()
		db.Close()
		return nil, fmt.Errorf("ошибка при подготовке запроса выбора: %w", err)
	}

	storage.preparedUpdate, err = db.Prepare("UPDATE links SET expiration_date = $1 WHERE original_url = $2")
	if err != nil {
		storage.closeAllStmts()
		db.Close()
		return nil, fmt.Errorf("ошибка при подготовке запроса обновления: %w", err)
	}

	storage.preparedInsert, err = storage.db.Prepare(`
        INSERT INTO links (short_url, original_url, expiration_date)
        VALUES ($1, $2, $3)
        ON CONFLICT (original_url) DO NOTHING
    `)
	if err != nil {
		storage.closeAllStmts()
		db.Close()
		return nil, fmt.Errorf("ошибка при подготовке запроса вставки: %w", err)
	}

	storage.preparedDelete, err = storage.db.Prepare("DELETE FROM links WHERE expiration_date < NOW()")
	if err != nil {
		storage.closeAllStmts()
		db.Close()
		return nil, fmt.Errorf("ошибка при подготовке запроса удаления: %w", err)
	}

	storage.preparedSelectOriginal, err = storage.db.Prepare("SELECT original_url FROM links WHERE short_url = $1")
	if err != nil {
		storage.closeAllStmts()
		db.Close()
		return nil, fmt.Errorf("ошибка при подготовке запроса: %w", err)
	}

	go storage.cleaningLoop()
	return storage, nil
}

// ShortenAndSaveUrl сокращает переданный Url и сохраняет сокращенный
func (d *DbStorage) ShortenAndSaveUrl(url string) (string, error) {

	var existingShort string
	err := d.preparedSelectShort.QueryRow(url).Scan(&existingShort)
	if err == nil {
		// запись найдена — обновляем expiration_date
		_, err = d.preparedUpdate.Exec(time.Now().Add(storageTime), url)
		if err != nil {
			return "", fmt.Errorf("ошибка при обновлении expiration_date: %w", err)
		}

		return existingShort, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("ошибка при проверке существующей ссылки: %w", err)
	}

	for attempts := 0; attempts < 25; attempts++ {
		shortUrl := short_link_generator.GenerateShortLink()
		_, err := d.preparedInsert.Exec(shortUrl, url, time.Now().Add(storageTime))
		if err != nil {
			return "", fmt.Errorf("ошибка при вставке ссылки: %w", err)
		}

		err = d.preparedSelectShort.QueryRow(url).Scan(&existingShort)
		if err == nil {
			return existingShort, nil
		}

	}

	return "", fmt.Errorf("не удалось сгенерировать уникальную короткую ссылку")
}

// GetUrl возвращает оригинальную ссылку по короткой, если она есть
func (d *DbStorage) GetUrl(shortUrl string) (string, error) {
	var url string
	err := d.preparedSelectOriginal.QueryRow(shortUrl).Scan(&url)

	if err != nil {
		return "", fmt.Errorf("не удалось найти короткую ссылку: %w", err)
	}

	_, err = d.preparedUpdate.Exec(time.Now().Add(storageTime), url)
	if err != nil {
		return url, fmt.Errorf("ошибка во время обновления expiration_date: %w", err)
	}
	return url, nil
}

// CleanOldRecords удаляет давно зарегестрированные ссылки
func (d *DbStorage) CleanOldRecords() error {

	_, err := d.preparedDelete.Exec()
	if err != nil {
		return fmt.Errorf("ошибка при удалении записей: %w", err)
	}
	return nil
}

// cleaningLoop удаляет давно зарегестрированные ссылки с определенным интервалом
func (d *DbStorage) cleaningLoop() {
	ticker := time.NewTicker(storageTime)
	defer ticker.Stop()
	for {
		select {
		case <-d.cleaningDone:
			d.cleaningFinished <- struct{}{}
			return
		case <-ticker.C:
			if err := d.CleanOldRecords(); err != nil {
				log.Printf("Ошибка при очистке устаревших записей: %v\n", err)
			}
		}
	}
}

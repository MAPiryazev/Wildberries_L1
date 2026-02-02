package storage

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// FileStorage - структура, для хранения базовой директории
type FileStorage struct {
	BaseDir string
}

// NewFileStorage - конструктор для FileStorage
func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{BaseDir: baseDir}
}

// SaveHTML сохраняет HTML страницы как downloads/<host>/index.html
func (s *FileStorage) SaveHTML(urlStr string, data []byte) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	dir := filepath.Join(s.BaseDir, u.Host)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("ошибка при создании директорий: %w", err)
	}

	path := filepath.Join(dir, "index.html")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("ошибка при сохранении HTML: %w", err)
	}
	return path, nil
}

// SaveResource сохраняет ресурс в структуру downloads/<host>/<resource_path>
func (s *FileStorage) SaveResource(urlStr string, data []byte) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Если путь пустой, даём дефолтное имя
	p := u.Path
	if p == "" || strings.HasSuffix(p, "/") {
		p = filepath.Join(p, "index")
	}

	localPath := filepath.Join(s.BaseDir, u.Host, p)

	// Создаём все родительские директории
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return "", fmt.Errorf("ошибка при создании директорий: %w", err)
	}

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return "", fmt.Errorf("ошибка при сохранении ресурса: %w", err)
	}

	return localPath, nil
}

// ExistsLocal проверяет, существует ли уже файл
func (s *FileStorage) ExistsLocal(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	localPath := filepath.Join(s.BaseDir, u.Host, u.Path)
	_, err = os.Stat(localPath)
	return err == nil
}

// DownloadToBytes скачивает любой URL в память
func DownloadToBytes(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при скачивании %s: %w", urlStr, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}
	return data, nil
}

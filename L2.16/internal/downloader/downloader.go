package downoader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gowget/internal/parser"
	"gowget/internal/storage"
)

// MaxDepth - константа для ограничения глубины рекурсии при парснге веб страниц
const MaxDepth int = 100

// Downloader - интерфейс, который имеет контракты для использования утилиты
type Downloader interface {
	DownloadPage(url string, limit int) error
}

type stockDownloader struct {
	Visited  map[string]bool
	MaxDepth int
	storage  *storage.FileStorage
}

// NewStockDownloader - конструктор для интерфейса downloader
func NewStockDownloader() *stockDownloader {
	return &stockDownloader{
		Visited:  make(map[string]bool),
		MaxDepth: MaxDepth,
		storage:  storage.NewFileStorage("downloads"),
	}
}

// resolveURL превращает относительный URL в абсолютный на основе baseURL
func resolveURL(baseURL, ref string) (string, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return "", err
	}
	if u.IsAbs() {
		return u.String(), nil
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(u).String(), nil
}

func (d *stockDownloader) DownloadPage(pageURL string, limit int) error {
	if limit < 0 {
		return nil
	}
	if d.Visited[pageURL] {
		return nil
	}
	d.Visited[pageURL] = true

	fmt.Printf("Скачиваю: %s (глубина: %d)\n", pageURL, limit)

	resp, err := http.Get(pageURL)
	if err != nil {
		return fmt.Errorf("ошибка при скачивании страницы %s: %w", pageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул статус %d для %s", resp.StatusCode, pageURL)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении тела страницы: %w", err)
	}

	if _, err := d.storage.SaveHTML(pageURL, bodyBytes); err != nil {
		return fmt.Errorf("ошибка при сохранении HTML: %w", err)
	}

	resources, err := parser.GetResources(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return fmt.Errorf("ошибка при парсинге HTML: %w", err)
	}

	// Собираем все ресурсы
	allRes := append(resources.CSS, resources.JS...)
	allRes = append(allRes, resources.Images...)
	allRes = append(allRes, resources.Media...)

	for _, res := range allRes {
		if res == "" {
			continue
		}

		fullURL, err := resolveURL(pageURL, res)
		if err != nil {
			fmt.Println("ошибка формирования полного URL ресурса:", res, err)
			continue
		}

		data, err := storage.DownloadToBytes(fullURL)
		if err != nil {
			fmt.Println("ошибка скачивания ресурса:", fullURL, err)
			continue
		}
		if _, err := d.storage.SaveResource(fullURL, data); err != nil {
			fmt.Println("ошибка сохранения ресурса:", fullURL, err)
		}
	}

	if limit > 0 {
		for _, link := range resources.Pages {
			if link == "" {
				continue
			}

			fullLink, err := resolveURL(pageURL, link)
			if err != nil {
				fmt.Println("ошибка формирования полного URL страницы:", link, err)
				continue
			}

			if err := d.DownloadPage(fullLink, limit-1); err != nil {
				fmt.Println("ошибка при скачивании подстраницы:", fullLink, err)
			}
		}
	}

	fmt.Printf("Страница сохранена: %s\n", pageURL)
	return nil
}

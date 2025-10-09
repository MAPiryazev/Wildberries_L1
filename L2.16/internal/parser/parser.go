package parser

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// PageLinks - структура, которая выступает хранилищем для контента веб страниц
type PageLinks struct {
	Pages  []string
	CSS    []string
	JS     []string
	Images []string
	Media  []string
	Other  []string
}

// GetResources парсит HTML-документ и достаёт все типы ссылок.
func GetResources(body io.Reader) (*PageLinks, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге страницы: %w", err)
	}
	if doc == nil {
		return nil, fmt.Errorf("html документ пустой")
	}

	links := &PageLinks{}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			GetLinks(n, links)
			GetCSS(n, links)
			GetJS(n, links)
			GetImages(n, links)
			GetMedia(n, links)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return links, nil
}

// GetLinks - функция для извлечения ссылок с web страниц
func GetLinks(n *html.Node, links *PageLinks) {
	if n.Data != "a" {
		return
	}
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			links.Pages = append(links.Pages, attr.Val)
		}
	}
}

// GetCSS - функция для извлечения CSS файлов с web страниц
func GetCSS(n *html.Node, links *PageLinks) {
	if n.Data != "link" {
		return
	}
	var isCSS bool
	for _, attr := range n.Attr {
		if attr.Key == "rel" && attr.Val == "stylesheet" {
			isCSS = true
			break
		}
	}
	if !isCSS {
		return
	}
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			links.CSS = append(links.CSS, attr.Val)
		}
	}
}

// GetJS - функция для извлечения JS файлов с web страниц
func GetJS(n *html.Node, links *PageLinks) {
	if n.Data != "script" {
		return
	}
	for _, attr := range n.Attr {
		if attr.Key == "src" {
			links.JS = append(links.JS, attr.Val)
		}
	}
}

// GetImages - функция для извлечения картинок с web страниц
func GetImages(n *html.Node, links *PageLinks) {
	if n.Data != "img" {
		return
	}
	for _, attr := range n.Attr {
		if attr.Key == "src" {
			links.Images = append(links.Images, attr.Val)
		}
	}
}

// GetMedia - функция для извлечения другого контента с web страниц
func GetMedia(n *html.Node, links *PageLinks) {
	switch n.Data {
	case "video", "audio", "source", "iframe":
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				links.Media = append(links.Media, attr.Val)
			}
		}
	}
}

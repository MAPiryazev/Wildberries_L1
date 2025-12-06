package templates

import (
	"html/template"
	"log"
)

var Templates *template.Template

// LoadTemplates загружает все html файлы из папки
func LoadTemplates(pattern string) {
	var err error
	Templates, err = template.ParseGlob(pattern)
	if err != nil {
		log.Fatalf("cannot parse templates: %v", err)
	}
}

package httpapi

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var embeddedStatic embed.FS

var assetsFS fs.FS

func init() {
	sub, err := fs.Sub(embeddedStatic, "web")
	if err != nil {
		panic(err)
	}
	assetsFS = sub
}

func serveIndex(w http.ResponseWriter, _ *http.Request) {
	data, err := fs.ReadFile(assetsFS, "index.html")
	if err != nil {
		http.Error(w, "index file not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

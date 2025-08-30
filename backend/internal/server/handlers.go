package server

import (
	"net/http"
	"path/filepath"
)

const (
	webDir = "web"
)

func BaseHandle(h http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(h, "method not supporting", http.StatusInternalServerError)
		return
	}
	requestedPath := r.URL.Path
	filePath := filepath.Join(webDir, requestedPath)

	_, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(h, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(h, r, filePath)

}

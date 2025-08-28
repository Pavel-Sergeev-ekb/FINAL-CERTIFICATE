package server

import (
	"database/sql"
	"net/http"
	"path/filepath"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/api"
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

func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.AddTaskHandler(db)(w, r)
		case http.MethodGet:
			api.GetTaskHandler(db)(w, r)
		case http.MethodPut:
			api.EditTaskHandler(db)(w, r)
		case http.MethodDelete:
			api.DeleteTaskHandler(db)(w, r)
		default:
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
		}
	}
}

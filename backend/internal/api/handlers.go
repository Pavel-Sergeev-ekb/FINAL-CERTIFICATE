package api

import (
	"net/http"
	"path/filepath"
	"time"
)

const (
	webDir = "web"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	//http.HandleFunc("/api/task", taskHandler)
	//http.HandleFunc("/api/tasks", tasksHandler)
	//http.HandleFunc("/api/task/done", doneTaskHandler)
}

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

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstartStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(TypeDay, nowStr)
		if err != nil {
			http.Error(w, ErrInvalidDate.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dstartStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(next))
}

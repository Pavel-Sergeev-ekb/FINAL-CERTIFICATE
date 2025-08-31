package api

import (
	"database/sql"
	"net/http"
)

func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddTaskHandler(db)(w, r)
		case http.MethodGet:
			GetTaskHandler(db)(w, r)
		case http.MethodPut:
			EditTaskHandler(db)(w, r)
		case http.MethodDelete:
			DeleteTaskHandler(db)(w, r)
		default:
			http.Error(w, `{"error":"method not supported"}`, http.StatusMethodNotAllowed)
		}
	}
}

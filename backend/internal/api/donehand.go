package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/services"
)

func DoneHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
			return
		}

		if db == nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "БД не подключена"})
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
			return
		}

		// Retrieve the task date and repeat rule
		var originalDateStr, repeatRule string
		err := db.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?", id).Scan(&originalDateStr, &repeatRule)
		if err == sql.ErrNoRows {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "Task not found"})
			return
		} else if err != nil {
			log.Printf("Error retrieving task: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
			return
		}

		// If repeat is empty → delete the task
		if repeatRule == "" {
			_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
			if err != nil {
				log.Printf("Error deleting task: %v", err)
				writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting task"})
				return
			}
			writeJson(w, http.StatusOK, map[string]string{})
			return
		}

		// Calculate the next execution date, passing `time.Now()` and the original date as a string
		nextDate, err := services.NextDate(time.Now(), originalDateStr, repeatRule)
		if err != nil {
			log.Printf("Error calculating the next date: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Error calculating the next date"})
			return
		}

		// Update the task date in the database
		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		if err != nil {
			log.Printf("Error updating task date: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Error updating task date"})
			return
		}

		writeJson(w, http.StatusOK, map[string]string{})
	}
}

package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/data"
)

func UpdateHandleTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task data.Task
		if db == nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "db not connection"})
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "JSON decoding error: " + err.Error()})
			return
		}
		if task.ID == 0 {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "task ID is required:"})
			return
		}
		if task.Title == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "task title not specified"})
			return
		}

		// Check if task exists
		var existingTask data.Task
		err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", task.ID).Scan(
			&existingTask.ID,
			&existingTask.Date,
			&existingTask.Title,
			&existingTask.Comment,
			&existingTask.Repeat,
		)
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		} else if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error task check: " + err.Error()})
			return
		}

		if _, err := time.Parse("20060102", task.Date); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid data format: " + err.Error()})
			return
		}
		if task.Repeat != "" && !isValidRepeat(task.Repeat) {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid repetition value"})
			return
		}

		// Update task in the database
		query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error update task: " + err.Error()})
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error getting number of updated records: " + err.Error()})
			return
		}

		if count == 0 {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		updatedTask := data.Task{}
		err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", task.ID).Scan(
			&updatedTask.ID,
			&updatedTask.Date,
			&updatedTask.Title,
			&updatedTask.Comment,
			&updatedTask.Repeat,
		)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error retrieving updated data"})
			return
		}

		// Формируем корректный JSON-ответ
		writeJson(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"task":   updatedTask,
		})
	}
}

package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/data"
)

func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Task ID is missing"})
			return
		}

		var task data.Task
		err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err == sql.ErrNoRows {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "Task not found"})
			return
		} else if err != nil {
			log.Printf("SQL Error: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving task"})
			return
		}

		// Проверяем, что ID корректно получен из БД
		if id == "" {
			log.Printf("Получен пустой ID задачи")
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Invalid task ID"})
			return
		}

		// Формируем ответ в правильном формате
		response := map[string]string{
			"id":      id,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}

		// Отправляем JSON-ответ
		writeJson(w, http.StatusOK, response)
	}
}

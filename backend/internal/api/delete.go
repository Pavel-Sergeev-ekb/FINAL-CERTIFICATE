package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
			return
		}

		if db == nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "БД не подключена"})
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "ID задачи не указан"})
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат ID задачи"})
			return
		}

		result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			log.Printf("Error deleting task: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении задачи"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
			return
		}

		log.Printf("Task with ID %d has been deleted", id)
		writeJson(w, http.StatusOK, map[string]interface{}{})
	}
}

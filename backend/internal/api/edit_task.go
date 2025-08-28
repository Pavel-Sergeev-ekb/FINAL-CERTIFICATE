package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/data"
)

func EditTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not supported"})
			return
		}
		if db == nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "БД не подключена"})
			return
		}

		// Decode JSON into a map[string]string to handle ID as a string
		var rawData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Ошибка декодирования JSON"})
			return
		}

		// Преобразуем ID в int64 с проверкой
		idStr, ok := rawData["id"].(string)
		if !ok {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат ID"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Ошибка преобразования ID"})
			return
		}

		// Формируем новую задачу с проверкой всех полей
		var newTask data.Task
		newTask.ID = id
		newTask.Date = rawData["date"].(string)
		newTask.Title = rawData["title"].(string)
		newTask.Comment = rawData["comment"].(string)
		newTask.Repeat = rawData["repeat"].(string)

		// Проверяем обязательные поля
		if newTask.ID == 0 {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "ID задачи обязателен"})
			return
		}
		if newTask.Title == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
			return
		}

		// Создаем новый JSON с корректным маршалином
		newBody, marshalErr := json.Marshal(newTask)
		if marshalErr != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка формирования JSON: " + marshalErr.Error()})
			return
		}

		// Заменяем тело запроса
		r.Body = io.NopCloser(bytes.NewReader(newBody))

		// Вызываем обработчик обновления
		UpdateHandleTask(db)(w, r)
	}
}

package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/data"
	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/services"
)

const formDay = "20060102"

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePostTask(db)(w, r)
		case http.MethodPut:
			UpdateHandleTask(db)(w, r)
		case http.MethodDelete:
			DeleteTaskHandler(db)(w, r)
		default:
			writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not supported"})
		}
	}
}

// HandleUpdateTask - processes task updates
func handlePostTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("new task addition request received")

		if db == nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "db not connection"})
			return
		}

		var task data.Task
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&task)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON format"})
			return
		}

		log.Printf("received task data: %+v", task)

		if task.Title == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "It is necessary to specify the task title"})
			return
		}

		nowStr := time.Now().Format(formDay)

		// Обработка даты
		if task.Date == "" || task.Date == "today" {
			task.Date = nowStr
			log.Printf("date not specified current set: %s", task.Date)
		}

		// Проверка формата даты
		_, err = time.Parse(formDay, task.Date)
		if err != nil {
			log.Printf("date format error: %s", task.Date)
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid date format"})
			return
		}

		// Обработка дат в прошлом
		if task.Date < nowStr {
			if task.Repeat == "" {
				task.Date = nowStr
				log.Printf("repeat not specified use now: %s", task.Date)
			} else {
				nextDate, err := services.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					log.Printf("error calculate repeat: %v", err)
					writeJson(w, http.StatusBadRequest, map[string]string{"error": "error calculate repeat"})
					return
				}
				task.Date = nextDate
				log.Printf("calculate next date: %s", task.Date)
			}
		}

		// Проверка значения повторения
		if task.Repeat != "" && !isValidRepeat(task.Repeat) {
			log.Printf("invalid repeat value: %s", task.Repeat)
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid repeat value"})
			return
		}

		// Вставка задачи в базу данных
		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			log.Printf("error adding task: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error adding task"})
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			log.Printf("error getting task ID: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error getting task ID"})
			return
		}

		log.Printf("Задача добавлена с ID: %d, дата: %s", id, task.Date)
		writeJson(w, http.StatusCreated, map[string]interface{}{"id": id})
	}
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func isValidRepeat(repeat string) bool {
	allowedPrefixes := []string{"d ", "y"}
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(repeat, prefix) {
			return true
		}
	}
	return false
}

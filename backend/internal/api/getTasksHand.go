package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty"` // string instead of int
	Date    string `json:"date"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if db == nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "БД не подключена"})
			return
		}

		// Retrieve the search parameter
		search := r.URL.Query().Get("search")
		var query string
		var args []interface{}
		const limit = 50 // Limit the number of tasks returned

		if search == "" {
			query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
			args = append(args, limit)
		} else {
			// Check if search is a date in the format "02.01.2006"
			parsedDate, err := time.Parse("02.01.2006", search)
			if err == nil {
				search = parsedDate.Format("20060102")
				query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?"
				args = append(args, search, limit)
			} else {
				query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
				args = append(args, "%"+search+"%", "%"+search+"%", limit)
			}
		}

		// Execute the SQL query
		rows, err := db.Query(query, args...)
		if err != nil {
			log.Printf("Error retrieving tasks: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Database query error"})
			return
		}
		defer rows.Close()

		// List of tasks
		var tasks []Task
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				log.Printf("Error reading row from database: %v", err)
				writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Data processing error"})
				return
			}
			tasks = append(tasks, task)
		}
		if tasks == nil {
			tasks = []Task{}
		}

		log.Printf("Sending %d tasks", len(tasks))

		// Отправка ответа
		writeJson(w, http.StatusOK, map[string]any{
			"tasks": tasks,
		})
	}
}

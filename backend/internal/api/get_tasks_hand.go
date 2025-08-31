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

const limit = 50
const typeDay = "20060102"

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if db == nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "db not connection"})
			return
		}

		search := r.URL.Query().Get("search")

		var query string
		var args []interface{}

		if search == "" {
			query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
			args = append(args, limit)
		} else {

			parsedDate, err := time.Parse("02.01.2006", search)
			if err == nil {
				search = parsedDate.Format(typeDay)
				query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?"
				args = append(args, search, limit)
			} else {
				query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
				args = append(args, "%"+search+"%", "%"+search+"%", limit)
			}
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Printf("error retrieving tasks: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "db query error"})
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				log.Printf("error reading row from database: %v", err)
				writeJson(w, http.StatusInternalServerError, map[string]string{"error": "data processing error"})
				return

			}
			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			log.Println(err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "database error"})
			return
		}

		if tasks == nil {
			tasks = []Task{}
		}

		log.Printf("Sending %d tasks", len(tasks))

		writeJson(w, http.StatusOK, map[string]any{
			"tasks": tasks,
		})
	}
}

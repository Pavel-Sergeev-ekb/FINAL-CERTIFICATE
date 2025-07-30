package data

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL, 
			comment TEXT, 
			date TIMESTAMP NOT NULL,
			repeat VARCHAR(50), 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_date ON scheduler(date);
`

func InitDB(db *sql.DB) error {

	dbPath := GetDBPath()

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec(schema)
		if err != nil {
			log.Fatal(err)
		}
	}
	return err
}

type Task struct {
	ID        int
	Date      string
	Title     string
	Comment   string
	Repeat    string
	createdAt time.Time
}

func FormatDate(date time.Time) string {
	return date.Format("20060102")
}

func GetDBPath() string {
	envPath := os.Getenv("TODO_DBFILE")

	if envPath == "" {
		return "scheduler.db"
	}

	dir := filepath.Dir(envPath)
	if dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatalf("not found dir for db: %s", dir)
		}
	}
	return envPath
}

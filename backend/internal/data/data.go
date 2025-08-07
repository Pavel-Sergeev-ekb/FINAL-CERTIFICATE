package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var idInt int64

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL, 
			comment TEXT, 
			date CHAR(8) NOT NULL DEFAULT '',
			repeat VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS  idx_scheduler_date ON scheduler(date);
`

func InitDB(dbFile string) error {

	_, err := os.Stat(dbFile)

	install := os.IsNotExist(err)

	var errOpen error

	db, errOpen := sql.Open("sqlite", dbFile)
	if errOpen != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	if install {
		if _, err = db.Exec(schema); err != nil {

			return fmt.Errorf("failed to init schema: %w", err)

		}

	}
	return nil
}

func DeleteTask(id string) error {
	if _, err := fmt.Sscan(id, &idInt); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, idInt)
	if err != nil {
		return fmt.Errorf("db delete error: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next, id string) error {
	if _, err := fmt.Sscan(id, &idInt); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}
	res, err := db.Exec(
		`UPDATE scheduler SET date = ? WHERE id = ?`, next, idInt,
	)
	if err != nil {
		return fmt.Errorf("db update date error: %w", err)
	}
	m, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if m == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
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

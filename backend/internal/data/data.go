package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL, 
			comment TEXT , 
			date CHAR(8) NOT NULL DEFAULT '',
			repeat VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS  idx_date ON scheduler(date);
`

func Init(dbFile string) (*sql.DB, error) {

	var db *sql.DB

	dbPath := GetDBPath()

	_, err := os.Stat(dbPath)

	install := os.IsNotExist(err)

	var errOpen error

	db, errOpen = sql.Open("sqlite", dbPath)
	if errOpen != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	if install {
		if _, err = db.Exec(schema); err != nil {

			return nil, fmt.Errorf("failed to init schema: %w", err)

		}

	}
	return db, nil
}

func DeleteTask(id string) error {

	var idInt int64
	var db *sql.DB

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

	var idInt int64
	var db *sql.DB

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

package data

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func NextDate(next string, id string) error {

	var db *sql.DB

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

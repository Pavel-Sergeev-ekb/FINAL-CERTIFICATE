package data

import (
	"database/sql"
	"fmt"
	"strings"
)

type Task struct {
	ID      int    `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

func CreateTask(task Task) (int64, error) {
	query := `
	INSERT INTO scheduler (data, title, comment, repeat)
	VALUES(?, ?, ?, ?)
	`
	res, err := db.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("db insert error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("db last insert id error: %w", err)
	}
	return id, nil
}

// Метод для получения задачи по ID
func GetTask(id string) (*Task, error) {
	var idInt int64
	if _, err := fmt.Sscan(id, &idInt); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	query := `
	SELECT id, date, title, comment, repeat
	FROM scheduler 
	WHERE id = ?
	LIMIT 1
	`
	row := db.QueryRow(query, idInt)

	var task Task
	var idOut int64

	switch err := row.Scan(
		&idOut,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	); err {
	case nil:
		//Если все хорошо - возвращаем задачу
		task.ID = int(idOut)
		return &task, nil
	case sql.ErrNoRows:
		// если запись не найдена
		return nil, fmt.Errorf("task not found")
	default:
		// для всех остальных ошибок
		return nil, fmt.Errorf("db scan error: %w", err)
	}
}

func UpdateTask(t *Task) error {
	if t.ID == 0 {
		return fmt.Errorf("invalid task ID")
	}

	query := `
	UPDATE scheduler
	SET date    = ?,
			title   = ?,
			comment = ?,
			repeat  = ?
			WHERE id = ?
			`

	res, err := db.Exec(query,
		t.Date,
		t.Title,
		t.Comment,
		t.Repeat,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("no task with id %d", t.ID)
	}
	return nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	baseQuery := `
 SELECT id, date, title, comment, repeat
 FROM scheduler
    `
	var args []interface{}
	where := ""
	// Попытка распознать DD.MM.YYYY
	if len(search) == len("02.01.2006") && strings.Count(search, ".") == 2 {
		// переводим в YYYYMMDD
		parts := strings.Split(search, ".")
		date := parts[2] + parts[1] + parts[0]
		where = " WHERE date = ?"
		args = append(args, date)
	} else if search != "" {
		where = " WHERE title LIKE ? OR comment LIKE ?"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	orderLimit := " ORDER BY date ASC LIMIT ?"
	args = append(args, limit)

	query := baseQuery + where + orderLimit

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.Tasks query error: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("db.Tasks scan error: %w", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db.Tasks rows error: %w", err)
	}

	// чтобы JSON с empty slice был [], а не null
	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}

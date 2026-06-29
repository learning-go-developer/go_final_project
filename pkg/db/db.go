package db

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const tableSchema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_sched_date ON scheduler(date);`

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Init(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if _, err = DB.Exec(tableSchema); err != nil {
		return err
	}

	return nil
}

func AddTask(t *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.ExecContext(context.Background(), query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetTask(id string) (*Task, error) {
	var t Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	if t.ID == "" {
		return nil, sql.ErrNoRows
	}
	return &t, nil
}

func UpdateTask(t *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return context.DeadlineExceeded
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func GetTasks() ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

package task

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type TaskEntity struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTasks(db *sql.DB, limit int) (interface{}, int, error) {
	var tasks []TaskEntity

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	for rows.Next() {
		var t TaskEntity
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(tasks) == 0 {
		tasks = []TaskEntity{}
	}

	return map[string][]TaskEntity{"tasks": tasks}, http.StatusOK, nil
}

func GetTaskByID(db *sql.DB, id string) ([]byte, int, error) {
	var t TaskEntity
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return []byte(`{"error": "Задача не найдена"}`), http.StatusNotFound, nil
		}
		return nil, http.StatusInternalServerError, err
	}

	response, err := json.Marshal(t)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return response, http.StatusOK, nil
}
